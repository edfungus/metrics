package registry

// Registry collects all the metrics to be stored
type Registry interface {
	Get(k Key) interface{}
	Filter(k Key) []interface{}
	Set(k Key, i interface{})
	Delete(k Key)
}

// Key is made up of key value pair combinations which can be filtered on later
type Key map[string]string
