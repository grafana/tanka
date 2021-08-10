package evalcache

import (
	"context"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GoogleStorageEvalCache is an evaluation cache that stores its data in Google Cloud Storage
// Cache paths with the gs:// scheme are routed to this cache
type GoogleStorageEvalCache struct {
	initMutex    sync.Once
	client       *storage.Client
	currentItems map[string]bool

	Bucket, Prefix string
}

func NewGoogleStorageEvalCache(url *url.URL) *GoogleStorageEvalCache {
	return &GoogleStorageEvalCache{
		Bucket: url.Host,
		Prefix: strings.Trim(url.Path, "/"),
	}
}

func (c *GoogleStorageEvalCache) cachePath(hash string) string {
	return filepath.Join(strings.Trim(c.Prefix, "/"), hash+".json")
}

func (c *GoogleStorageEvalCache) Get(hash string) (string, error) {
	ctx := context.Background()

	var err error
	c.initMutex.Do(func() {
		log.Printf("Initializing the GCS cache. Bucket: %s, Prefix: %s", c.Bucket, c.Prefix)
		if c.client, err = storage.NewClient(ctx); err != nil {
			return
		}

		query := &storage.Query{Prefix: c.Prefix}

		bkt := c.client.Bucket(c.Bucket)
		it := bkt.Objects(ctx, query)
		c.currentItems = map[string]bool{}
		for {
			var attrs *storage.ObjectAttrs
			attrs, err = it.Next()
			if err == iterator.Done {
				err = nil
				break
			}
			if err != nil {
				return
			}
			c.currentItems[attrs.Name] = true
		}
	})
	if err != nil {
		return "", err
	}

	cachePath := c.cachePath(hash)
	if _, ok := c.currentItems[cachePath]; ok {
		reader, err := c.client.Bucket(c.Bucket).Object(cachePath).NewReader(ctx)
		if err != nil {
			return "", err
		}
		bytes, err := io.ReadAll(reader)
		return string(bytes), err
	}

	return "", nil
}

func (c *GoogleStorageEvalCache) Store(hash, content string) error {
	ctx := context.Background()

	cachePath := c.cachePath(hash)
	writer := c.client.Bucket(c.Bucket).Object(cachePath).NewWriter(ctx)
	if _, err := io.WriteString(writer, content); err != nil {
		return err
	}
	c.currentItems[cachePath] = true

	return writer.Close()
}
