package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-baa/cache/lru"
)

const (
	// MemoryLimit default memory size limit, 128mb
	// 128  1 << 7
	// 1024 1 << 10
	MemoryLimit int64 = 1 << 27
)

// Memory implement a memory cache adapter for cacher
type Memory struct {
	Name       string
	Prefix     string
	bytes      int64
	bytesLimit int64
	mu         sync.RWMutex
	store      *lru.Cache
}

// Exist return true if value cached by given key
func (c *Memory) Exist(key string) bool {
	item := c.get(c.Prefix + key)
	if item != nil {
		return true
	}
	return false
}

// Get returns value by given key
func (c *Memory) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item := c.get(c.Prefix + key)
	if item != nil {
		return item.Val
	}
	return nil
}

func (c *Memory) get(key string) *Item {
	v, ok := c.store.Get(c.Prefix + key)
	if !ok {
		return nil
	}
	item, err := v.(ItemBinary).Item()
	if err != nil {
		return nil
	}
	if item.Expired() {
		c.Delete(c.Prefix + key)
		return nil
	}
	return item
}

// Set cache value by given key
func (c *Memory) Set(key string, v interface{}, ttl int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := NewItem(v, ttl)
	b, err := item.Bytes()
	if err != nil {
		return err
	}
	// if overwrite bytes count will error
	// so, delete first if exist
	if c.Exist(key) {
		c.store.Remove(c.Prefix + key)
	}
	l := int64(len(b))
	c.gc(l)
	c.store.Add(c.Prefix+key, b)
	c.bytes += l
	return nil
}

// Incr increases cached int-type value by given key as a counter
// if key not exist, before increase set value with zero
func (c *Memory) Incr(key string) (interface{}, error) {
	item := c.get(key)
	if item == nil {
		item = NewItem(0, 0)
	}
	err := item.Incr()
	if err != nil {
		return nil, err
	}
	ttl := int64((item.Expiration - time.Now().UnixNano()) / 1e9)
	if ttl < 0 {
		return nil, fmt.Errorf("cache expired")
	}
	err = c.Set(key, item.Val, ttl)
	if err != nil {
		return nil, err
	}
	return item.Val, nil
}

// Decr decreases cached int-type value by given key as a counter
// if key not exist, return errors
func (c *Memory) Decr(key string) (interface{}, error) {
	item := c.get(key)
	if item == nil {
		return nil, fmt.Errorf("cache key not exists")
	}
	err := item.Decr()
	if err != nil {
		return nil, err
	}
	ttl := int64((item.Expiration - time.Now().UnixNano()) / 1e9)
	if ttl < 0 {
		return nil, fmt.Errorf("cache expired")
	}
	err = c.Set(key, item.Val, ttl)
	if err != nil {
		return nil, err
	}
	return item.Val, nil
}

// Delete delete cached data by given key
func (c *Memory) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Remove(c.Prefix + key)
	return nil
}

// Flush flush cacher
func (c *Memory) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var l int
	for {
		l = c.store.Len()
		if l > 0 {
			c.store.RemoveOldest()
		} else {
			break
		}
	}
	c.bytes = 0

	return nil
}

// Start new a cacher and start service
func (c *Memory) Start(o Options) error {
	c.Name = o.Name
	c.Prefix = o.Prefix
	c.bytesLimit = MemoryLimit
	if o.Config != nil {
		if v, ok := o.Config["bytesLimit"].(int64); ok {
			c.bytesLimit = v
		}
	}

	if c.store == nil {
		c.store = lru.New(0)
		c.store.OnEvicted = func(key lru.Key, value interface{}) {
			c.bytes -= int64(len(value.(ItemBinary)))
		}
	}

	return nil
}

// gc release memory for storage new item
// if free bytes can store item returns
// remove items until bytes less than bytesLimit - len(item) * 64
func (c *Memory) gc(size int64) {
	if c.bytes+size < c.bytesLimit {
		return
	}

	var l int
	for c.bytes > c.bytesLimit-size*64 {
		l = c.store.Len()
		if l > 0 {
			c.store.RemoveOldest()
		} else {
			break
		}
	}
	c.bytes = 0
}

func init() {
	Register("memory", &Memory{})
}
