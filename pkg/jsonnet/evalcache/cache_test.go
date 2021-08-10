package evalcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCache(t *testing.T) {
	cases := []struct {
		name     string
		opts     CacheOpts
		expected EvalCache
		err      string
	}{
		{
			name: "no scheme",
			opts: CacheOpts{CachePath: "/tmp/test-folder"},
			err:  "unhandled caching URL scheme: /tmp/test-folder",
		},
		{
			name: "bad scheme",
			opts: CacheOpts{CachePath: "test:///tmp/test-folder"},
			err:  "unhandled caching URL scheme: test:///tmp/test-folder",
		},
		{
			name:     "local path",
			opts:     CacheOpts{CachePath: "file:///tmp/test-folder"},
			expected: &FileEvalCache{Directory: "/tmp/test-folder"},
		},
		{
			name:     "local path is trimmed",
			opts:     CacheOpts{CachePath: "file:///tmp/test-folder/"},
			expected: &FileEvalCache{Directory: "/tmp/test-folder"},
		},
		{
			name:     "gcs bucket",
			opts:     CacheOpts{CachePath: "gs://test-bucket/test-folder/nested"},
			expected: &GoogleStorageEvalCache{Bucket: "test-bucket", Prefix: "test-folder/nested"},
		},
		{
			name:     "gcs prefix is trimmed",
			opts:     CacheOpts{CachePath: "gs://test-bucket/test-folder/nested/"},
			expected: &GoogleStorageEvalCache{Bucket: "test-bucket", Prefix: "test-folder/nested"},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for k := range cacheMap {
				delete(cacheMap, k)
			}

			result, err := GetCache(tc.opts)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}
