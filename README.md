# registry

Store your values (metrics, counters, etc) in a data structure that can retrieve your exact entry or a list of entries that matches a given key where a key is a key value pair.

### Goals

Be able to write/read metrics with keys that can be filtered on. That way multiple metrics entries can have common keys for easy retrieval and analytical metric calculations. In the future, this should be ale to performantly log metrics in a time interval. 

Also, thread safety will come later. Will add locks!

### Implementations

* simple.go has a straight forward naive implementation of registry 
* better.go has a better(?) implementation of registry
* cached.go has the better implmentation with a cache
