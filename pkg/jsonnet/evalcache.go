package jsonnet

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileEvalCache is an evaluation cache that stores its data on the local filesystem
type FileEvalCache struct {
	Directory string
}

func NewFileEvalCache(cachePath string) *FileEvalCache {
	return &FileEvalCache{
		Directory: cachePath,
	}
}

func (c *FileEvalCache) cachePath(hash string) (string, error) {
	return filepath.Abs(filepath.Join(c.Directory, hash+".json"))
}

func (c *FileEvalCache) Get(hash string) (string, error) {
	cachePath, err := c.cachePath(hash)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(cachePath); err == nil {
		bytes, err := ioutil.ReadFile(cachePath)
		return string(bytes), err
	} else if !os.IsNotExist(err) {
		return "", err
	}
	return "", nil
}

func (c *FileEvalCache) Store(hash, content string) error {
	if err := os.MkdirAll(c.Directory, os.ModePerm); err != nil {
		return err
	}

	cachePath, err := c.cachePath(hash)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, []byte(content), 0644)
}
