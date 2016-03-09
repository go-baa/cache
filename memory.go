package cache

import (
	"sync"
	"time"
)

// Memory implement a memory cache adapter for cacher
type Memory struct {
	Name     string
	store    map[string]*item
	lock     sync.RWMutex
	interval int64
}

type item struct {
	v       interface{}
	created int64
	ttl     int64
}

// Expired check item has expired
func (t *item) Expired() bool {
	return t.ttl > 0 && (time.Now().Unix()-t.created) >= t.ttl
}

// Exist check key is exist
func (c *Memory) Exist(key string) bool {
	_, ok := c.store[key]
	if ok {
		return !c.store[key].Expired()
	}
	return false
}

// Get returns value for given key
func (c *Memory) Get(key string) interface{} {
	if c.Exist(key) {
		return c.store[key].v
	}
	return nil
}

// Set set value for given key
func (c *Memory) Set(key string, v interface{}, ttl int64) error {
	c.store[key] = &item{
		v:       v,
		created: time.Now().Unix(),
		ttl:     ttl,
	}
	return nil
}

// Delete delete the key
func (c *Memory) Delete(key string) error {
	delete(c.store, key)
	return nil
}

// Flush flush cacher
func (c *Memory) Flush() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key := range c.store {
		delete(c.store, key)
	}
	return nil
}

// Start new a cacher and start service
func (c *Memory) Start(o Options) error {
	if c.store == nil {
		c.store = make(map[string]*item)
	}
	c.interval = o.Interval

	go c.gc()

	return nil
}

func (c *Memory) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.interval < 1 {
		return
	}

	if c.store != nil {
		for key := range c.store {
			if c.store[key].Expired() {
				delete(c.store, key)
			}
		}
	}

	time.AfterFunc(time.Duration(c.interval)*time.Second, func() {
		c.gc()
	})
}

func init() {
	Register("memory", &Memory{})
}
