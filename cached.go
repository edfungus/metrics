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
	entries, err := c.getCache.Get(k)
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
	return entry.value
}
func (c *CachedRegistry) Filter(k Key) []Entry {
	entries, err := c.filterCache.Get(k)
	if err == nil {
		return toEntryArray(entries)
	}
	entries = c.registry.Filter(k)
	cacheItemRemoved, cacheKeyRemoved, cacheValueRemoved := c.filterCache.UpdateWithKey(k, entries)
	if cacheItemRemoved {
		removeFilterCacheKey(cacheKeyRemoved, cacheValueRemoved.(hashEntries))
	}
	return toEntryArray(entries)
}

func (c *CachedRegistry) Set(k Key, i interface{}) {

}
func (c *CachedRegistry) Delete(k Key) {

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
