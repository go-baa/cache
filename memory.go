package cache

import (
	"sync"
	"time"

	"github.com/go-baa/cache/lru"
)

const (
	// MemoryLimit default memory size limit
	// 128  1 << 7
	// 1024 1 << 10
	// 128M 1 << 27
	MemoryLimit int64 = 1 << 27
)

// Memory implement a memory cache adapter for cacher
type Memory struct {
	Name       string
	bytes      int64
	bytesLimit int64
	mu         sync.RWMutex
	store      *lru.Cache
}

// Exist check key is exist
func (c *Memory) Exist(key string) bool {
	item := c.get(key)
	if item != nil {
		return true
	}
	return false
}

// Get returns value for given key
func (c *Memory) Get(key string) interface{} {
	item := c.get(key)
	if item != nil {
		return item.Val
	}
	return nil
}

func (c *Memory) get(key string) *Item {
	v, ok := c.store.Get(key)
	if !ok {
		return nil
	}
	item := new(Item)
	err := DecodeGob(v.([]byte), item)
	if err != nil {
		return nil
	}
	if item.Expired() {
		c.Delete(key)
		return nil
	}
	return item
}

// Set set value for given key
func (c *Memory) Set(key string, v interface{}, ttl int64) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item := &Item{
		Val:     v,
		Created: time.Now().Unix(),
		TTL:     ttl,
	}
	b, err := EncodeGob(item)
	if err != nil {
		return err
	}
	l := int64(len(b))
	c.gc(l)
	c.store.Add(key, b)
	c.bytes += l
	return nil
}

// Delete delete the key
func (c *Memory) Delete(key string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.store.Remove(key)
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
	c.bytesLimit = MemoryLimit
	if o.Config != nil {
		if v, ok := o.Config["bytesLimit"].(int64); ok {
			c.bytesLimit = v
		}
	}

	if c.store == nil {
		c.store = lru.New(0)
		c.store.OnEvicted = func(key lru.Key, value interface{}) {
			c.bytes -= int64(len(value.([]byte)))
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
