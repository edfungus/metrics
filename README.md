# registry

Store your values (metrics, counters, etc) in a data structure that can retrieve your exact entry or a list of entries that matches a given key where a key is a key value pair.

### Goals

Be able to write/read metrics with keys that can be filtered on. That way multiple metrics entries can have common keys for easy retrieval and analytical metric calculations. In the future, this should be ale to performantly log metrics in a time interval. 

### Implementations

* simple.go has a straight forward dumb implementation of registry 
