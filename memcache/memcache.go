package memcache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-baa/cache"
)

// Memcache implement a memcache cache adapter for cacher
type Memcache struct {
	Name   string
	handle *memcache.Client
}

// Exist check key is exist
func (c *Memcache) Exist(key string) bool {
	val := c.Get(key)
	if val != nil {
		return true
	}
	return false
}

// Get returns value for given key
func (c *Memcache) Get(key string) interface{} {
	v, err := c.handle.Get(key)
	if err != nil {
		return nil
	}
	item, err := cache.ItemBinary(v.Value).Item()
	if err != nil || item == nil {
		return nil
	}
	return item.Val
}

// Set set value for given key
func (c *Memcache) Set(key string, v interface{}, ttl int64) error {
	item := cache.NewItem(v, ttl)
	b, err := item.Bytes()
	if err != nil {
		return err
	}
	return c.handle.Set(&memcache.Item{Key: key, Value: []byte(b), Expiration: int32(ttl)})
}

// Delete delete the key
func (c *Memcache) Delete(key string) error {
	return c.handle.Delete(key)
}

// Flush flush cacher
func (c *Memcache) Flush() error {
	return c.handle.FlushAll()
}

// Start new a cacher and start service
func (c *Memcache) Start(o cache.Options) error {
	c.Name = o.Name
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
	err := c.handle.Set(&memcache.Item{Key: "foo", Value: []byte("bar")})
	if err != nil {
		return fmt.Errorf("memcache connect err: %s", err)
	}
	return nil
}

func init() {
	cache.Register("memcache", &Memcache{})
}
