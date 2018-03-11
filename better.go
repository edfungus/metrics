package registry

/*

	Seeing that using keys and their values are usually repeated, we can take
	advantage of this by using maps to perhaps speed up the key search. The performance
	here is a little bit harder to calculate and is also based upon uniqueness of keys. This
	should be faster than simpleRegistry, but by how much? Probably need some perf tests.

*/

// BetterRegistry is a better implmentation of registry .. but is it really better?
type BetterRegistry struct {
	registry map[string]values
}
type values map[string]entries
type entries []*Entry

func NewBetterRegistry() *BetterRegistry {
	return &BetterRegistry{
		registry: map[string]values{},
	}
}

func (b *BetterRegistry) Get(k Key) interface{} {
	entry := getEntryWithKey(b.registry, k)
	if entry == nil {
		return nil
	}
	return entry.Value
}

func (b *BetterRegistry) Filter(k Key) []Entry {
	entriesWithKey := getEntriesContainKey(b.registry, k)
	entries := []Entry{}
	for _, entry := range entriesWithKey {
		entries = append(entries, *entry)
	}
	return entries
}

func (b *BetterRegistry) Set(k Key, i interface{}) {
	entryWithKey := getEntryWithKey(b.registry, k)
	if entryWithKey != nil {
		entryWithKey.Value = i
	} else {
		entryWithKey = &Entry{
			Key:   k,
			Value: i,
		}
		addEntry(b.registry, entryWithKey)

	}
}

func (b *BetterRegistry) Delete(k Key) {
	entriesWithKey := getEntryWithKey(b.registry, k)
	if entriesWithKey == nil {
		return
	}
	removeEntry(b.registry, entriesWithKey)
}

func addEntry(r map[string]values, e *Entry) {
	for key, value := range e.Key {
		_, ok := r[key]
		if !ok {
			r[key] = values{}
		}
		r[key][value] = append(r[key][value], e)
	}
}

func removeEntry(r map[string]values, e *Entry) {
	for key, value := range e.Key {
		if _, ok := r[key]; !ok {
			return
		}
		entries, ok := r[key][value]
		if !ok {
			return
		}
		for i, entry := range entries {
			if entry == e {
				r[key][value] = append(r[key][value][:i], r[key][value][i+1:]...)
				break
			}
		}
		if len(r[key][value]) == 0 {
			delete(r[key], value)
		}
		if len(r[key]) == 0 {
			delete(r, key)
		}
	}
}

// getEntryWithKey returns the Entry that has the exact Key
func getEntryWithKey(r map[string]values, k Key) *Entry {
	entriesWithKey := getEntriesContainKey(r, k)
	for _, entry := range entriesWithKey {
		if len(entry.Key) == len(k) {
			return entry
		}
	}
	return nil
}

// getEntriesContainKey returns all entries that contain all the key value pairs in Key
func getEntriesContainKey(r map[string]values, k Key) entries {
	entriesThatContainKey := []entries{}
	for key, value := range k {
		entriesWithAKey := getEntriesWithAKey(r, key, value)
		if len(entriesWithAKey) == 0 {
			return entries{}
		}
		entriesThatContainKey = append(entriesThatContainKey, entriesWithAKey)
	}
	return findUnionOfEntires(entriesThatContainKey)
}

// findUnionOfEntires returns a list of entries that are commom between all arrays of entries
func findUnionOfEntires(es []entries) entries {
	m := map[*Entry]int{}
	for _, e := range es {
		for _, entry := range e {
			m[entry]++
		}
	}
	entries := entries{}
	for e, count := range m {
		if count == len(es) {
			entries = append(entries, e)
		}
	}
	return entries
}

// getEntriesWithAKey returns all entires that has this key value pair in the Key
func getEntriesWithAKey(r map[string]values, k string, v string) entries {
	values, ok := r[k]
	if !ok {
		return entries{}
	}
	entriesWithKey, ok := values[v]
	if !ok {
		return entries{}
	}
	return entriesWithKey
}
