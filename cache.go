// Package cache providers a cache management for baa.
package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

// Cacher a cache management for baa
type Cacher interface {
	// Exist return true if value cached by given key
	Exist(key string) bool
	// Get returns value to o by given key
	Get(key string, o interface{}) error
	// Set cache value by given key
	Set(key string, v interface{}, ttl int64) error
	// Incr increases cached int-type value by given key as a counter
	// if key not exist, before increase set value with zero
	Incr(key string) (int64, error)
	// Decr decreases cached int-type value by given key as a counter
	// if key not exist, return errors
	Decr(key string) (int64, error)
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
		t.Val = int64(t.Val.(int)) + 1
	case int8:
		t.Val = int64(t.Val.(int8)) + 1
	case int16:
		t.Val = int64(t.Val.(int16)) + 1
	case int32:
		t.Val = int64(t.Val.(int32)) + 1
	case int64:
		t.Val = int64(t.Val.(int64)) + 1
	case uint:
		t.Val = int64(t.Val.(uint)) + 1
	case uint8:
		t.Val = int64(t.Val.(uint8)) + 1
	case uint16:
		t.Val = int64(t.Val.(uint16)) + 1
	case uint32:
		t.Val = int64(t.Val.(uint32)) + 1
	case uint64:
		t.Val = int64(t.Val.(uint64)) + 1
	default:
		return fmt.Errorf("item value is not int-type")
	}
	return nil
}

// Decr decreases given value
func (t *Item) Decr() error {
	switch t.Val.(type) {
	case int:
		t.Val = int64(t.Val.(int)) - 1
	case int8:
		t.Val = int64(t.Val.(int8)) - 1
	case int16:
		t.Val = int64(t.Val.(int16)) - 1
	case int32:
		t.Val = int64(t.Val.(int32)) - 1
	case int64:
		t.Val = int64(t.Val.(int64)) - 1
	case uint:
		if t.Val.(uint) > 0 {
			t.Val = int64(t.Val.(uint)) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint8:
		if t.Val.(uint8) > 0 {
			t.Val = int64(t.Val.(uint8)) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint16:
		if t.Val.(uint16) > 0 {
			t.Val = int64(t.Val.(uint16)) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint32:
		if t.Val.(uint32) > 0 {
			t.Val = int64(t.Val.(uint32)) - 1
		} else {
			return fmt.Errorf("item value is less than 0")
		}
	case uint64:
		if t.Val.(uint64) > 0 {
			t.Val = int64(t.Val.(uint64)) - 1
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

// SimpleType check value type is simple type or not
func SimpleType(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	case bool:
		return true
	default:
		return false
	}
}

// SimpleValue return value to output with type convert
func SimpleValue(v []byte, o interface{}) bool {
	switch o.(type) {
	case *string:
		*o.(*string) = string(v)
	case *bool:
		*o.(*bool), _ = strconv.ParseBool(string(v))
	case *int:
		t, _ := strconv.ParseInt(string(v), 10, 64)
		*o.(*int) = int(t)
	case *int8:
		t, _ := strconv.ParseInt(string(v), 10, 64)
		*o.(*int8) = int8(t)
	case *int16:
		t, _ := strconv.ParseInt(string(v), 10, 64)
		*o.(*int16) = int16(t)
	case *int32:
		t, _ := strconv.ParseInt(string(v), 10, 64)
		*o.(*int32) = int32(t)
	case *int64:
		*o.(*int64), _ = strconv.ParseInt(string(v), 10, 64)
	case *uint:
		t, _ := strconv.ParseUint(string(v), 10, 64)
		*o.(*uint) = uint(t)
	case *uint8:
		t, _ := strconv.ParseUint(string(v), 10, 64)
		*o.(*uint8) = uint8(t)
	case *uint16:
		t, _ := strconv.ParseUint(string(v), 10, 64)
		*o.(*uint16) = uint16(t)
	case *uint32:
		t, _ := strconv.ParseUint(string(v), 10, 64)
		*o.(*uint32) = uint32(t)
	case *uint64:
		*o.(*uint64), _ = strconv.ParseUint(string(v), 10, 64)
	case *float32:
		t, _ := strconv.ParseFloat(string(v), 64)
		*o.(*float32) = float32(t)
	case *float64:
		*o.(*float64), _ = strconv.ParseFloat(string(v), 64)
	default:
		return false
	}
	return true
}

func init() {
	gob.Register(time.Time{})
	gob.Register(&Item{})
}
