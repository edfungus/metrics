package registry

import "errors"

// SimpleRegistry is a quick and dirty dumb implementation of the resgistry ... using just for benchmarking!
type SimpleRegistry struct {
	registry []*simpleRegistryMetric
}

type simpleRegistryMetric struct {
	key   map[string]string
	value interface{}
}

var (
	keyNotFound = errors.New("Key not found")
)

// NewSimpleRegistry returns a simple implementation of registry
func NewSimpleRegistry() *SimpleRegistry {
	return &SimpleRegistry{
		registry: []*simpleRegistryMetric{},
	}
}

// Get returns the metric that matches the key exactly
func (r *SimpleRegistry) Get(k Key) interface{} {
	v, _, err := getMetric(r.registry, k)
	if err != nil {
		return nil
	}
	return v.value
}

// Filter returns a list of metrics that matches the key
func (r *SimpleRegistry) Filter(k Key) []Entry {
	ks := []Entry{}
	for _, m := range r.registry {
		if isSubset(k, m.key) {
			ks = append(ks, Entry{
				Key:   m.key,
				Value: m.value,
			})
		}
	}
	return ks
}

// Set replaces or creates new entry with key and value
func (r *SimpleRegistry) Set(k Key, i interface{}) {
	m, _, err := getMetric(r.registry, k)
	if err == keyNotFound {
		m := &simpleRegistryMetric{
			key:   k,
			value: i,
		}
		r.registry = append(r.registry, m)
	}
	if m != nil {
		m.value = i
	}
}

// Delete removes an entry from the registry
func (r *SimpleRegistry) Delete(k Key) {
	_, i, err := getMetric(r.registry, k)
	if err == keyNotFound {
		return
	}
	r.registry = append(r.registry[:i], r.registry[i+1:]...)
}

func getMetric(ms []*simpleRegistryMetric, k Key) (*simpleRegistryMetric, int, error) {
	for i, m := range ms {
		if isEquals(m.key, k) {
			return m, i, nil
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
