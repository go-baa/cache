package memcache

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-baa/cache"
)

// Memcache implement a memcache cache adapter for cacher
type Memcache struct {
	Name   string
	Prefix string
	handle *memcache.Client
}

// Exist return true if value cached by given key
func (c *Memcache) Exist(key string) bool {
	_, err := c.handle.Get(c.Prefix + key)
	if err == nil {
		return true
	}
	return false
}

// Get returns value by given key
func (c *Memcache) Get(key string, out interface{}) error {
	v, err := c.handle.Get(c.Prefix + key)
	if err != nil {
		return err
	}

	if cache.SimpleValue(v.Value, out) {
		return nil
	}

	item, err := cache.ItemBinary(v.Value).Item()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(out)
	if rv.IsNil() {
		return errors.New("cache: out is nil")
	}
	if rv.Kind() != reflect.Ptr {
		return errors.New("cache: out must be a pointer")
	}
	for rv.Kind() == reflect.Ptr {
		if !rv.Elem().IsValid() && rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	if !rv.CanSet() {
		return errors.New("cache: out cannot set value")
	}
	if rv.Type() != reflect.TypeOf(item.Val) {
		return fmt.Errorf("cache: out is different type with stored value %v, %v", rv.Type(), reflect.TypeOf(item.Val))
	}
	rv.Set(reflect.ValueOf(item.Val))
	return nil
}

// Set cache value by given key
func (c *Memcache) Set(key string, v interface{}, ttl int64) error {
	var t []byte
	if !cache.SimpleType(v) {
		item := cache.NewItem(v, ttl)
		b, err := item.Bytes()
		if err != nil {
			return err
		}
		t = []byte(b)
	} else {
		t = []byte(fmt.Sprintf("%v", v))
	}
	return c.handle.Set(&memcache.Item{Key: c.Prefix + key, Value: t, Expiration: int32(ttl)})
}

// Incr increases cached int-type value by given key as a counter
// if key not exist, before increase set value with zero
func (c *Memcache) Incr(key string) (int64, error) {
	v, err := c.handle.Increment(c.Prefix+key, 1)
	return int64(v), err
}

// Decr decreases cached int-type value by given key as a counter
// if key not exist, return errors
func (c *Memcache) Decr(key string) (int64, error) {
	v, err := c.handle.Decrement(c.Prefix+key, 1)
	return int64(v), err
}

// Delete delete cached data by given key
func (c *Memcache) Delete(key string) error {
	return c.handle.Delete(c.Prefix + key)
}

// Flush flush cacher
func (c *Memcache) Flush() error {
	return c.handle.FlushAll()
}

// Start new a cacher and start service
func (c *Memcache) Start(o cache.Options) error {
	c.Name = o.Name
	c.Prefix = o.Prefix
	var host, port string
	if val, ok := o.Config["host"]; ok {
		host = val.(string)
	} else {
		host = "127.0.0.1"
	}
	if val, ok := o.Config["port"]; ok {
		port = val.(string)
	} else {
		port = "11211"
	}

	c.handle = memcache.New(host + ":" + port)
	err := c.handle.Set(&memcache.Item{Key: c.Prefix + "foo", Value: []byte("bar")})
	if err != nil {
		return fmt.Errorf("memcache connect err: %s", err)
	}
	return nil
}

func init() {
	cache.Register("memcache", &Memcache{})
}
