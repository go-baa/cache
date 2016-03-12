// Package cache providers a cache management for baa.
package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// Cacher a cache management for baa
type Cacher interface {
	// Exist check key is exist
	Exist(key string) bool
	// Get returns value for given key
	Get(key string) interface{}
	// Set set value for given key
	Set(key string, v interface{}, ttl int64) error
	// Delete delete the key
	Delete(key string) error
	// Flush flush cacher
	Flush() error
	// Start new a cacher and start service
	Start(Options) error
}

// Item cache storage item
type Item struct {
	Val     interface{}
	Created int64
	TTL     int64
}

// Options cache options
type Options struct {
	Name     string                 // cache name
	Adapter  string                 // adapter
	Config   map[string]interface{} // config for adapter
	Interval int64                  // gc interval second
}

var adapters = make(map[string]Cacher)

// New create a Cacher
func New(o Options) Cacher {
	if o.Name == "" {
		o.Name = "_DEFAULT_"
	}
	if o.Adapter == "" {
		panic("cache.New: cannot use empty adapter")
	}
	c, err := NewCacher(o.Adapter, o)
	if err != nil {
		panic("cache.New: " + err.Error())
	}
	return c
}

// NewCacher creates and returns a new cacher by given adapter name and configuration.
// It panics when given adapter isn't registered and starts GC automatically.
func NewCacher(name string, o Options) (Cacher, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", name)
	}
	return adapter, adapter.Start(o)
}

// Register registers a adapter
func Register(name string, adapter Cacher) {
	if adapter == nil {
		panic("cache.Register: cannot register adapter with nil value")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("cache.Register: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}

// Expired check item has expired
func (t *Item) Expired() bool {
	return t.TTL > 0 && (time.Now().Unix()-t.Created) >= t.TTL
}

// EncodeGob encode item use gob for storage
func EncodeGob(item *Item) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(item)
	return buf.Bytes(), err
}

// DecodeGob decode item use gob for storage
func DecodeGob(data []byte, out *Item) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(&out)
}

func init() {
	gob.Register(time.Time{})
	gob.Register(&Item{})
}
