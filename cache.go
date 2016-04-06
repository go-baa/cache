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
	// Exist return true if value cached by given key
	Exist(key string) bool
	// Get returns value by given key
	Get(key string) interface{}
	// Set cache value by given key
	Set(key string, v interface{}, ttl int64) error
	// Incr increases cached int-type value by given key as a counter
	// if key not exist, before increase set value with zero
	Incr(key string) (interface{}, error)
	// Decr decreases cached int-type value by given key as a counter
	// if key not exist, return errors
	Decr(key string) (interface{}, error)
	// Delete delete cached data by given key
	Delete(key string) error
	// Flush flush cacher
	Flush() error
	// Start new a cacher and start service
	Start(Options) error
}

// Item cache storage item
type Item struct {
	Val        interface{} // real object value
	TTL        int64       // cache life time
	Expiration int64       // expired time
}

// ItemBinary cache item encoded data
type ItemBinary []byte

// Options cache options
type Options struct {
	Name    string                 // cache name
	Adapter string                 // adapter
	Prefix  string                 // cache key prefix
	Config  map[string]interface{} // config for adapter
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

// NewItem create a cache item
func NewItem(val interface{}, ttl int64) *Item {
	item := &Item{Val: val, TTL: ttl}
	if ttl > 0 {
		item.Expiration = time.Now().Add(time.Duration(ttl) * time.Second).UnixNano()
	}
	return item
}

// Expired check item has expired
func (t *Item) Expired() bool {
	return t.TTL > 0 && time.Now().UnixNano() >= t.Expiration
}

// Incr increases given value
func (t *Item) Incr() error {
	switch t.Val.(type) {
	case int:
		t.Val = t.Val.(int) + 1
	case int32:
		t.Val = t.Val.(int32) + 1
	case int64:
		t.Val = t.Val.(int64) + 1
	case uint:
		t.Val = t.Val.(uint) + 1
	case uint32:
		t.Val = t.Val.(uint32) + 1
	case uint64:
		t.Val = t.Val.(uint64) + 1
	default:
		return fmt.Errorf("item value is not int-type")
	}
	return nil
}

// Decr decreases given value
func (t *Item) Decr() error {
	switch t.Val.(type) {
	case int:
		t.Val = t.Val.(int) - 1
	case int32:
		t.Val = t.Val.(int32) - 1
	case int64:
		t.Val = t.Val.(int64) - 1
	case uint:
		if t.Val.(uint) > 0 {
			t.Val = t.Val.(uint) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint32:
		if t.Val.(uint32) > 0 {
			t.Val = t.Val.(uint32) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint64:
		if t.Val.(uint64) > 0 {
			t.Val = t.Val.(uint64) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	default:
		return fmt.Errorf("item value is not int-type")
	}
	return nil
}

// Bytes encode item use gob for storage
func (t *Item) Bytes() (ItemBinary, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(t)
	return buf.Bytes(), err
}

// Item decode bytes data to cache item use gob
func (t ItemBinary) Item() (*Item, error) {
	buf := bytes.NewBuffer(t)
	item := new(Item)
	err := gob.NewDecoder(buf).Decode(&item)
	return item, err
}

func init() {
	gob.Register(time.Time{})
	gob.Register(&Item{})
}
