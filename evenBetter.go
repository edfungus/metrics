package registry

/*

	As I was writing CachedRegistry, I realized that I should just rewrite Better Registry
	with new specs and use that in CachedRegistry so things will be cleaner. This is
	conceptually BetterRegistry with things moved around with better organization

*/

type EvenBetterRegistry struct {
	registry map[string]map[string]hashEntries
	// keys    *hashKey // where string is key and value respectively
}

// // hashKey is essentially a Key with one value
// type hashKey struct {
// 	key   string
// 	value string
// 	// cacheKeys []string
// }

// hashEntry is essentially a Entry with an array which includes all the hashes this Entry is in
type hashEntry struct {
	keys            map[string]string
	value           interface{}
	getCacheKey     string   // A formed hash key (could be empty) for which this hashEntry is referred to in the cache
	filterCacheKeys []string // List of form hash keys (string) for which this hashEntry is referred to in the cache
}

func NewEvenBetterRegistry() *EvenBetterRegistry {
	return &EvenBetterRegistry{
		registry: map[string]map[string]hashEntries{},
		// keys:     map[string]map[string]*hashKey{},
	}
}

func (r *EvenBetterRegistry) Get(k Key) (*hashEntry, error) {
	hashEntries := r.getHashEntriesForKey(k)
	entries := findUnionOfHashEntries(hashEntries)
	var hashEntry *hashEntry
	for _, entry := range entries {
		if len(entry.keys) == len(k) {
			hashEntry = entry
		}
	}
	if hashEntry == nil {
		return nil, keyNotFound
	}
	hashEntry.getCacheKey = toHashString(k)
	return hashEntry, nil
}
func (r *EvenBetterRegistry) Filter(k Key) hashEntries {
	hashEntries := r.getHashEntriesForKey(k)
	return findUnionOfHashEntries(hashEntries)
}
func (r *EvenBetterRegistry) Set(k Key, i interface{}) {

}
func (r *EvenBetterRegistry) Delete(k Key) {

}

// getHashEntriesForKey returns all hashKeys that contain Key (which means the Key can contain MORE than the given Key)
func (r *EvenBetterRegistry) getHashEntriesForKey(k Key) map[string]hashEntries {
	hashKeys := map[string]hashEntries{}
	for key, value := range k {
		_, ok := r.registry[key]
		if !ok {
			continue
		}
		hashEntry, ok := r.registry[key][value]
		if !ok {
			continue
		}
		hashKeys[key] = hashEntry
	}
	return hashKeys
}

func findUnionOfHashEntries(m map[string]hashEntries) hashEntries {
	counter := map[*hashEntry]int{}
	for _, entries := range m {
		for _, entry := range entries {
			counter[entry]++
		}
	}
	entries := hashEntries{}
	for entry, count := range counter {
		if count == len(m) {
			entries = append(entries, entry)
		}
	}
	return entries
}
