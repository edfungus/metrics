package registry

import "errors"

/*

	This is probably the most straight forward implmentation. In terms of peformance,
	it is pretty much linear time for all of the functions. It is going to be
	O(n + avg key size).

	Also it is not thread safe. Just need to add some locks for that

*/

// SimpleRegistry is a quick and dirty naive implementation of the resgistry
type SimpleRegistry struct {
	registry []*Entry
}

var (
	keyNotFound = errors.New("Key not found")
)

// NewSimpleRegistry returns a simple implementation of registry
func NewSimpleRegistry() *SimpleRegistry {
	return &SimpleRegistry{
		registry: []*Entry{},
	}
}

// Get returns the metric that matches the key exactly
func (r *SimpleRegistry) Get(k Key) interface{} {
	e, _, err := getEntry(r.registry, k)
	if err != nil {
		return nil
	}
	return e.Value
}

// Filter returns a list of metrics that matches the key
func (r *SimpleRegistry) Filter(k Key) []Entry {
	ks := []Entry{}
	for _, e := range r.registry {
		if isSubset(k, e.Key) {
			ks = append(ks, *e)
		}
	}
	return ks
}

// Set replaces or creates new entry with key and value
func (r *SimpleRegistry) Set(k Key, i interface{}) {
	e, _, err := getEntry(r.registry, k)
	if err == keyNotFound {
		e := &Entry{
			Key:   k,
			Value: i,
		}
		r.registry = append(r.registry, e)
	}
	if e != nil {
		e.Value = i
	}
}

// Delete removes an entry from the registry
func (r *SimpleRegistry) Delete(k Key) {
	_, i, err := getEntry(r.registry, k)
	if err == keyNotFound {
		return
	}
	r.registry = append(r.registry[:i], r.registry[i+1:]...)
}

func getEntry(es []*Entry, k Key) (*Entry, int, error) {
	for i, e := range es {
		if isEquals(e.Key, k) {
			return e, i, nil
		}
	}
	return nil, 0, keyNotFound
}

func isEquals(k1 Key, k2 Key) bool {
	if len(k1) != len(k2) {
		return false
	}
	return isSubset(k1, k2)
}

// isSubset checks whether of not all of the Keys in k is in the set
func isSubset(k Key, set Key) bool {
	for key := range k {
		sv, ok := set[key]
		if !ok {
			return false
		}
		if sv != k[key] {
			return false
		}
	}
	return true
}
