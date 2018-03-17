package registry

import (
	"reflect"
)

/*

	As I was writing CachedRegistry, I realized that I should just rewrite Better Registry
	with new specs and use that in CachedRegistry so things will be cleaner. This is
	conceptually BetterRegistry with things moved around with better organization

*/

type EvenBetterRegistry struct {
	registry map[string]map[string]hashEntries
}

// hashEntry is essentially a Entry with an array which includes all the hashes this Entry is in
type hashEntry struct {
	keys            Key
	value           interface{}
	getCacheKey     string   // A formed hash key (could be empty) for which this hashEntry is referred to in the cache
	filterCacheKeys []string // List of form hash keys (string) for which this hashEntry is referred to in the cache
}

func NewEvenBetterRegistry() *EvenBetterRegistry {
	return &EvenBetterRegistry{
		registry: map[string]map[string]hashEntries{},
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
	entry, err := r.Get(k)
	if err != nil {
		newEntry := &hashEntry{
			keys:            k,
			value:           i,
			getCacheKey:     "",
			filterCacheKeys: []string{},
		}
		r.addHashEntry(newEntry)
		return
	}
	entry.value = i
}

func (r *EvenBetterRegistry) Delete(k Key) *hashEntry {
	var entry *hashEntry
	for key, value := range k {
		e := r.removeEntryFromAKey(key, value, k)
		if e == nil {
			return nil
		}
		entry = e
	}
	return entry
}

func (r *EvenBetterRegistry) removeEntryFromAKey(key string, value string, completeKey Key) *hashEntry {
	values, ok := r.registry[key]
	if !ok {
		return nil
	}
	hashEntries, ok := values[value]
	if !ok {
		return nil
	}
	hashEntries, entry := removeFromHashEntries(hashEntries, completeKey)
	if entry == nil {
		return nil
	}
	// Some clean up before leaving
	if len(hashEntries) == 0 {
		delete(values, value)
	}
	if len(values) == 0 {
		delete(r.registry, key)
	}
	return entry
}

func removeFromHashEntries(entries hashEntries, k Key) (hashEntries, *hashEntry) {
	for i, e := range entries {
		if reflect.DeepEqual(e.keys, k) {
			entries = append(entries[:i], entries[i+1:]...)
			return entries, e
		}
	}
	return entries, nil
}

// getHashEntriesForKey returns all hashKeys that contain Key (which means the Key can contain MORE than the given Key)
func (r *EvenBetterRegistry) getHashEntriesForKey(k Key) map[string]hashEntries {
	hashKeys := map[string]hashEntries{}
	for key, value := range k {
		_, ok := r.registry[key]
		if !ok {
			hashKeys[key] = hashEntries{}
			continue
		}
		hashEntry, ok := r.registry[key][value]
		if !ok {
			hashKeys[key] = hashEntries{}
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

func (r *EvenBetterRegistry) addHashEntry(e *hashEntry) {
	for key, value := range e.keys {
		_, ok := r.registry[key]
		if !ok {
			r.registry[key] = map[string]hashEntries{}
		}
		r.registry[key][value] = append(r.registry[key][value], e)
	}
}
