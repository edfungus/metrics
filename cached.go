package registry

/*

	This is essentially BetterRegistry with a cache. However, we couldn't directly use
	BetterRegistry because we needed the Key and Entry to be able to relate back to the
	cache so it can be updated

*/

const (
	cacheSize = 10
)

type CachedRegistry struct {
	registry    *EvenBetterRegistry
	getCache    Cache // caches gets. Should only have one Entry per cache key
	filterCache Cache // cache filters. Will have a list of Entries that satisfies the Key
}

type hashEntries []*hashEntry

func NewCacheRegistry() *CachedRegistry {
	return &CachedRegistry{
		registry:    NewEvenBetterRegistry(),
		getCache:    NewSimpleCache(cacheSize),
		filterCache: NewSimpleCache(cacheSize),
	}
}

func (c *CachedRegistry) Get(k Key) interface{} {
	entries, err := c.getCache.GetWithKey(k)
	if err == nil {
		return entries.(hashEntries)[0].value
	}
	entry, err := c.registry.Get(k)
	if err != nil {
		return nil
	}
	cacheItemRemoved, _, cacheValueRemoved := c.getCache.UpdateWithHash(entry.getCacheKey, hashEntries{entry})
	if cacheItemRemoved {
		removeGetCacheKey(cacheValueRemoved.(hashEntries))
	}
	addGetCacheKey(entry, k)
	return entry.value
}
func (c *CachedRegistry) Filter(k Key) []Entry {
	entries, err := c.filterCache.GetWithKey(k)
	if err == nil {
		return toEntryArray(entries)
	}
	entries = c.registry.Filter(k)
	cacheItemRemoved, cacheKeyRemoved, cacheValueRemoved := c.filterCache.UpdateWithKey(k, entries)
	if cacheItemRemoved {
		removeFilterCacheKey(cacheKeyRemoved, cacheValueRemoved.(hashEntries))
	}
	addFilterCacheKey(entries.(hashEntries), k)
	return toEntryArray(entries)
}

func (c *CachedRegistry) Set(k Key, i interface{}) {
	c.registry.Set(k, i)
}

func (c *CachedRegistry) Delete(k Key) {
	entry := c.registry.Delete(k)
	if entry == nil {
		return
	}
	// Clean up entry from caches
	c.getCache.RemoveWithHash(entry.getCacheKey)
	for _, hashString := range entry.filterCacheKeys {
		entries, err := c.filterCache.GetWithHash(hashString)
		if err != nil {
			// shouldn't get here, this means entry cache list and cache are out of sync
			continue
		}
		entries, _ = removeFromHashEntries(entries.(hashEntries), k)
		if len(entries.(hashEntries)) == 0 {
			c.filterCache.RemoveWithHash(hashString)
		} else {
			c.filterCache.UpdateWithHash(hashString, entries)
		}
	}
}

/*
	We want to update the hashEntry that has been removed from the cache by deleteing the getCacheKey
		which holds the key for the hashEntry in the cache
	deleteEntries should really only have one value
*/
func removeGetCacheKey(deleteEntries hashEntries) {
	for _, entry := range deleteEntries {
		entry.getCacheKey = ""
	}
}

func addGetCacheKey(entry *hashEntry, k Key) {
	entry.getCacheKey = toHashString(k)
}

/*
	Given the cache is filled, UpdateWithHash() will give us the cache key and value.
	The value given will be an array of pointers to hashEntry(s) which we must remove the now deleted cache key
		from the filterCacheKeys. Therefore this updates the hashEntry(s) that it has been removed from the cache for this key
*/
func removeFilterCacheKey(deletedKey string, deleteEntries hashEntries) {
	for _, entry := range deleteEntries {
		for i, hash := range entry.filterCacheKeys {
			if hash == deletedKey {
				entry.filterCacheKeys = append(entry.filterCacheKeys[:i], entry.filterCacheKeys[i+1:]...)
				break
			}
		}
	}
}

func addFilterCacheKey(entries hashEntries, k Key) {
	cacheKey := toHashString(k)
	for _, entry := range entries {
		entry.filterCacheKeys = append(entry.filterCacheKeys, cacheKey)
	}
}

// toEntryArray converts hashEntries (in interface or hashEntries type form) to []Entry
func toEntryArray(i interface{}) []Entry {
	entries := []Entry{}
	for _, entry := range i.(hashEntries) {
		entries = append(entries, Entry{
			Key:   entry.keys,
			Value: entry.value,
		})
	}
	return entries
}
