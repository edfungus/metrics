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
	return []Entry{}
}
func (c *CachedRegistry) Set(k Key, i interface{}) {

}
func (c *CachedRegistry) Delete(k Key) {

}

func removeGetCacheKey(entries hashEntries) {
	for _, entry := range entries {
		entry.getCacheKey = ""
	}
}
