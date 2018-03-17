package registry

import (
	"encoding/base64"
	"sort"
)

const (
	hashKeyValueDelimiter = ":"
	hashKeyDelimiter      = ","
)

type Cache interface {
	GetWithKey(k Key) (interface{}, error)
	GetWithHash(h string) (interface{}, error)
	UpdateWithKey(k Key, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{})
	UpdateWithHash(h string, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{})
	RemoveWithHash(h string)
}

// SimpleCache will store commomly requested Keys and the associated Entries that have those keys
// The Entries will contain at least the Key BUT could have more than the specified keys
type SimpleCache struct {
	cache      map[string]interface{} // Caches something based on hashed string
	recentKeys []string               // In sync with the cache to keep track of cache entry recency
	maxSize    int                    // Max size of cache
}

func NewSimpleCache(maxSize int) *SimpleCache {
	return &SimpleCache{
		cache:      map[string]interface{}{},
		recentKeys: []string{}, // TODO: Could use make() to define length later
		maxSize:    maxSize,
	}
}

// GetWithKey gets a value from the cache that matches complete key
func (c *SimpleCache) GetWithKey(k Key) (interface{}, error) {
	hashString := toHashString(k)
	return c.GetWithHash(hashString)

}

// GetWithHash gets a value from the cache that matches the hashstring
func (c *SimpleCache) GetWithHash(h string) (interface{}, error) {
	value := c.cache[h]
	if value != nil {
		return value, nil
	}
	return nil, keyNotFound
}

// UpdateWithKey adds a value to the cache. If cache size will be exceed, the oldest value is removed and also returned
func (c *SimpleCache) UpdateWithKey(k Key, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{}) {
	return c.UpdateWithHash(toHashString(k), i)
}

// UpdateWithHash adds a value to the cache. If cache size will be exceed, the oldest value is removed and also returned
func (c *SimpleCache) UpdateWithHash(h string, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{}) {
	if _, ok := c.cache[h]; ok {
		c.RemoveWithHash(h)
	}
	c.cache[h] = i
	c.addToRecentKeys(h)
	return c.checkCacheSize()
}

func (c *SimpleCache) RemoveWithHash(h string) {
	c.removeFromRecentKeys(h)
	delete(c.cache, h)
}

// checkCacheSize checks if size has met maxSize and if so, remove oldest cache item
// NOTE: This only removes ONE value ... in case of uncertain overage, use a loop
func (c *SimpleCache) checkCacheSize() (exceeded bool, cacheKeyRemoved string, cacheValueRemoved interface{}) {
	if len(c.cache) <= c.maxSize {
		return false, "", nil
	}
	cacheKeyRemoved = c.recentKeys[0]
	cacheValueRemoved = c.cache[cacheKeyRemoved]
	c.RemoveWithHash(cacheKeyRemoved)
	return true, cacheKeyRemoved, cacheValueRemoved
}

func (c *SimpleCache) removeFromRecentKeys(h string) {
	for i, hash := range c.recentKeys {
		if hash == h {
			c.recentKeys = append(c.recentKeys[:i], c.recentKeys[i+1:]...)
		}
	}
}

func (c *SimpleCache) addToRecentKeys(h string) {
	c.removeFromRecentKeys(h)
	c.recentKeys = append(c.recentKeys, h)
}

func toHashString(k Key) string {
	hash := ""
	sortedKey := sortedKeys(k)
	for _, key := range sortedKey {
		encodedKey := base64.StdEncoding.EncodeToString([]byte(key))
		encodedValue := base64.StdEncoding.EncodeToString([]byte(k[key]))
		hash = hash + encodedKey + hashKeyValueDelimiter + encodedValue + hashKeyDelimiter
	}
	return hash
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, len(m), len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
