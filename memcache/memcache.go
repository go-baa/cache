package memcache

import (
	"fmt"
	"time"

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
	val := c.Get(c.Prefix + key)
	if val != nil {
		return true
	}
	return false
}

// Get returns value by given key
func (c *Memcache) Get(key string) interface{} {
	v, err := c.handle.Get(c.Prefix + key)
	if err != nil {
		return nil
	}
	item, err := cache.ItemBinary(v.Value).Item()
	if err != nil || item == nil {
		return nil
	}
	return item.Val
}

// Set cache value by given key
func (c *Memcache) Set(key string, v interface{}, ttl int64) error {
	item := cache.NewItem(v, ttl)
	b, err := item.Bytes()
	if err != nil {
		return err
	}
	return c.handle.Set(&memcache.Item{Key: c.Prefix + key, Value: []byte(b), Expiration: int32(ttl)})
}

// Incr increases cached int-type value by given key as a counter
// if key not exist, before increase set value with zero
func (c *Memcache) Incr(key string) (interface{}, error) {
	v, err := c.handle.Get(c.Prefix + key)
	if err != nil {
		if err.Error() == memcache.ErrCacheMiss.Error() {
			var v interface{}
			v = 1
			err = c.Set(key, v, 0)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
		return nil, err
	}
	item, err := cache.ItemBinary(v.Value).Item()
	if err != nil || item == nil {
		return nil, err
	}
	err = item.Incr()
	if err != nil {
		return nil, err
	}
	b, err := item.Bytes()
	if err != nil {
		return nil, err
	}
	ttl := int64((item.Expiration - time.Now().UnixNano()) / 1e9)
	if ttl < 0 {
		return nil, fmt.Errorf("cache expired")
	}
	err = c.handle.Set(&memcache.Item{Key: c.Prefix + key, Value: []byte(b), Expiration: int32(ttl)})
	return item.Val, nil
}

// Decr decreases cached int-type value by given key as a counter
// if key not exist, return errors
func (c *Memcache) Decr(key string) (interface{}, error) {
	v, err := c.handle.Get(c.Prefix + key)
	if err != nil {
		return nil, err
	}
	item, err := cache.ItemBinary(v.Value).Item()
	if err != nil || item == nil {
		return nil, err
	}
	err = item.Decr()
	if err != nil {
		return nil, err
	}
	b, err := item.Bytes()
	if err != nil {
		return nil, err
	}
	ttl := int64((item.Expiration - time.Now().UnixNano()) / 1e9)
	if ttl < 0 {
		return nil, fmt.Errorf("cache expired")
	}
	err = c.handle.Set(&memcache.Item{Key: c.Prefix + key, Value: []byte(b), Expiration: int32(ttl)})
	return item.Val, nil
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
