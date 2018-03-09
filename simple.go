package registry

// SimpleRegistry is a quick and dirty dumb implementation of the resgistry ... using just for benchmarking!
type SimpleRegistry struct {
	registry []simpleRegistryMetric
}

type simpleRegistryMetric struct {
	key   map[string]string
	value interface{}
}

// NewSimpleRegistry returns a simple implementation of registry
func NewSimpleRegistry() *SimpleRegistry {
	return &SimpleRegistry{
		registry: []simpleRegistryMetric{},
	}
}

// Get returns the metric that matches the key exactly
func (r *SimpleRegistry) Get(k Key) interface{} {
	// for _, m := range r.registry {

	// }
	return nil
}

// Filter returns a list of metrics that matches the key
func (r *SimpleRegistry) Filter(k Key) []interface{} {
	return []interface{}{}
}

// Set
func (r *SimpleRegistry) Set(k Key, i interface{}) error {
	return nil
}

func (r *SimpleRegistry) Delete(k Key) {

}

func isEquals(k1 Key, k2 Key) bool {
	if len(k1) != len(k2) {
		return false
	}
	for key := range k1 {
		v2, ok := k2[key]
		if !ok {
			return false
		}
		if v2 != k1[key] {
			return false
		}
	}
	return true
}
