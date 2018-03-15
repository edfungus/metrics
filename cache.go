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
	Get(k Key) (interface{}, error)
	UpdateWithKey(k Key, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{})
	UpdateWithHash(h string, i interface{}) (cacheItemRemoved bool, cacheKeyRemoved string, cacheValueRemoved interface{})
}

// SimpleCache will store commomly requested Keys and the associated Entries that have those keys
// The Entries will contain at least the Key BUT could have more than the specified keys
type SimpleCache struct {
	cache      map[string]interface{}
	recentKeys []string
	maxSize    int
}

func NewSimpleCache(maxSize int) *SimpleCache {
	return &SimpleCache{
		cache:      map[string]interface{}{},
		recentKeys: []string{}, // TODO: Could use make() to define length later
		maxSize:    maxSize,
	}
}

// Get get a value from the cache
func (c *SimpleCache) Get(k Key) (interface{}, error) {
	hashString := toHashString(k)
	value := c.cache[hashString]
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
		c.removeWithHash(h)
	}
	c.cache[h] = i
	c.addToRecentKeys(h)
	return c.checkCacheSize()
}

// checkCacheSize checks if size has met maxSize and if so, remove oldest cache item
// NOTE: This only removes ONE value ... in case of uncertain overage, use a loop
func (c *SimpleCache) checkCacheSize() (exceeded bool, cacheKeyRemoved string, cacheValueRemoved interface{}) {
	if len(c.cache) <= c.maxSize {
		return false, "", nil
	}
	cacheKeyRemoved = c.recentKeys[0]
	cacheValueRemoved = c.cache[cacheKeyRemoved]
	c.removeWithHash(cacheKeyRemoved)
	return true, cacheKeyRemoved, cacheValueRemoved
}

func (c *SimpleCache) removeWithHash(h string) {
	c.removeFromRecentKeys(h)
	delete(c.cache, h)
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
