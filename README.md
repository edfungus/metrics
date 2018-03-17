# registry

Store your values (metrics, counters, etc) in a data structure that can retrieve your exact entry or a list of entries that matches a given key where a key is a key value pair.

### Intro

Be able to write/read metrics with keys that can be filtered on. That way multiple metrics entries can have common keys for easy retrieval and analytical metric calculations. An example use case to be able to performantly log metrics in a time interval. 

```go
type Registry interface {
	Get(k Key) interface{}
	Filter(k Key) []Entry
	Set(k Key, i interface{})
	Delete(k Key)
}

```

The best implementation to use is `cachedRegistry` (`r := NewCacheRegistry(cacheSize)`) which combines the better registry implementation with a cache!

Also, thread safety will come later. Will add locks!

### Implementations

* `simple.go` has a straight forward naive implementation of registry 
* `better.go` has a better(?) implementation of registry
* `evenBetter.go` has an even better(!) implementation of registry. Really, its the same as `better.go` just rearranged and works with `cachedRegistry`
* `cached.go` has the even better implmentation with a cache
