package evalcache

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GoogleStorageEvalCache is an evaluation cache that stores its data in Google Cloud Storage
// Cache paths with the gs:// scheme are routed to this cache
type GoogleStorageEvalCache struct {
	initMutex    sync.Once
	client       *minio.Client
	currentItems map[string]bool

	Endpoint, Bucket, Prefix string
}

func NewGoogleStorageEvalCache(url *url.URL) *GoogleStorageEvalCache {
	return &GoogleStorageEvalCache{
		Endpoint: "https://storage.googleapis.com",
		Bucket:   url.Host,
		Prefix:   strings.Trim(url.Path, "/"),
	}
}

func (c *GoogleStorageEvalCache) cachePath(hash string) string {
	return filepath.Join(strings.Trim(c.Prefix, "/"), hash+".json")
}

func (c *GoogleStorageEvalCache) Get(hash string) (string, error) {
	ctx := context.Background()

	var err error
	c.initMutex.Do(func() {
		log.Printf("Initializing the S3 cache. Endpont: %s, Bucket: %s, Prefix: %s", c.Endpoint, c.Bucket, c.Prefix)
		if c.client, err = minio.New(c.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("GCS_ACCESS_KEY_ID"), os.Getenv("GCS_SECRET_ACCESS_KEY"), ""),
			Secure: true,
		}); err != nil {
			return
		}

		c.currentItems = map[string]bool{}
		objectCh := c.client.ListObjects(ctx, c.Bucket, minio.ListObjectsOptions{
			Prefix:    c.Prefix,
			Recursive: false,
		})
		for object := range objectCh {
			if err = object.Err; err != nil {
				return
			}
			c.currentItems[object.Key] = true
		}
	})
	if err != nil {
		return "", err
	}

	cachePath := c.cachePath(hash)
	if _, ok := c.currentItems[cachePath]; ok {
		object, err := c.client.GetObject(context.Background(), c.Bucket, cachePath, minio.GetObjectOptions{})
		if err != nil {
			return "", err
		}
		bytes, err := io.ReadAll(object)
		return string(bytes), err
	}

	return "", nil
}

func (c *GoogleStorageEvalCache) Store(hash, content string) error {
	ctx := context.Background()

	buffer := &bytes.Buffer{}
	if _, err := io.WriteString(buffer, content); err != nil {
		return err
	}

	cachePath := c.cachePath(hash)
	_, err := c.client.PutObject(ctx, c.Bucket, cachePath, buffer, int64(buffer.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	c.currentItems[cachePath] = true

	return nil
}
