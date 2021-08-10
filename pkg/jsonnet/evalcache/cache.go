package evalcache

import (
	"errors"
	"net/url"
	"regexp"
	"sync"
)

// CacheOpts represents configurable parameters for Jsonnet evaluation caching
type CacheOpts struct {
	CachePath        string
	CachePathRegexes []*regexp.Regexp
}

// PathMatches determines if a given path is matched by any of the configured path regexes
// If no path regexes are defined, all paths are matched
func (opts CacheOpts) PathMatches(path string) bool {
	for _, regex := range opts.CachePathRegexes {
		if regex.MatchString(path) {
			return true
		}
	}
	return len(opts.CachePathRegexes) == 0
}

// EvalCache represents a means to store and retrieve long-running evaluation results
// The key is a hash that represents the jsonnet code and its dependencies
type EvalCache interface {
	Get(hash string) (string, error)
	Store(hash, content string) error
}

// Contains caches that have been initialized.
// This reduces required processing and
//  ensures that all caches are initialized only once.
var cacheMutex sync.Mutex
var cacheMap = map[string]EvalCache{}

// GetCache gets or create a cache instance for the given cache options
// Only one cache for cache path can exist, as some caches require initialization
func GetCache(opts CacheOpts) (EvalCache, error) {
	if opts.CachePath == "" {
		return nil, nil
	}

	url, err := url.ParseRequestURI(opts.CachePath)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if v, ok := cacheMap[opts.CachePath]; ok {
		return v, nil
	}

	var cache EvalCache
	switch url.Scheme {
	case "file":
		cache = NewFileEvalCache(url)
	case "gs":
		cache = NewGoogleStorageEvalCache(url)
	default:
		return nil, errors.New("unhandled caching URL scheme: " + opts.CachePath)
	}

	cacheMap[opts.CachePath] = cache
	return cache, nil
}
