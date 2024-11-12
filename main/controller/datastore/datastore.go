package datastore

import (
	"sync"
	"time"
)

type datainput struct {
	value  interface{}
	update time.Time
}

var datastore map[interface{}]interface{}

var mu sync.Mutex

func init() {
	datastore = make(map[interface{}]interface{})
}

func Add(t, d interface{}) {
	mu.Lock()
	defer mu.Unlock()
	str := datainput{
		value:  d,
		update: time.Now(),
	}
	datastore[t] = str
}

func Read(t interface{}) interface{} {
	mu.Lock()
	defer mu.Unlock()
	tmp, ok := datastore[t].(datainput)
	if !ok {
		return nil
	}
	return tmp.value
}

func ReadTime(t interface{}) time.Time {
	mu.Lock()
	defer mu.Unlock()

	tmp, ok := datastore[t].(datainput)
	if !ok {
		return time.Time{}
	}
	return tmp.update
}
