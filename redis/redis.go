package redis

import (
	"fmt"
	"time"

	"github.com/go-baa/cache"
	"gopkg.in/redis.v3"
)

// Redis implement a redis cache adapter for cacher
type Redis struct {
	Name   string
	Prefix string
	handle *redis.Client
}

// Exist return true if value cached by given key
func (c *Redis) Exist(key string) bool {
	ok, err := c.handle.Exists(c.Prefix + key).Result()
	if err == nil && ok {
		return true
	}
	return false
}

// Get returns value by given key
func (c *Redis) Get(key string, out interface{}) error {
	v, err := c.handle.Get(c.Prefix + key).Bytes()
	if err != nil {
		return err
	}

	if cache.SimpleValue(v, out) {
		return nil
	}

	item, err := cache.ItemBinary(v).Item()
	if err != nil {
		return err
	}
	return item.Decode(out)
}

// Set cache value by given key
func (c *Redis) Set(key string, v interface{}, ttl int64) error {
	if !cache.SimpleType(v) {
		item := cache.NewItem(v, ttl)
		b, err := item.Encode()
		if err != nil {
			return err
		}
		v = []byte(b)
	}
	return c.handle.Set(c.Prefix+key, v, time.Second*time.Duration(ttl)).Err()
}

// Incr increases cached int-type value by given key as a counter
// if key not exist, before increase set value with zero
func (c *Redis) Incr(key string) (int64, error) {
	t := c.handle.Incr(c.Prefix + key)
	if t.Err() != nil {
		return 0, t.Err()
	}
	return t.Val(), nil
}

// Decr decreases cached int-type value by given key as a counter
// if key not exist, return errors
func (c *Redis) Decr(key string) (int64, error) {
	t := c.handle.Decr(c.Prefix + key)
	if t.Err() != nil {
		return 0, t.Err()
	}
	return t.Val(), nil
}

// Delete delete cached data by given key
func (c *Redis) Delete(key string) error {
	return c.handle.Del(c.Prefix + key).Err()
}

// Flush flush cacher
func (c *Redis) Flush() error {
	return c.handle.FlushDb().Err()
}

// Start new a cacher and start service
func (c *Redis) Start(o cache.Options) error {
	c.Name = o.Name
	c.Prefix = o.Prefix
	var host, port, pass string
	var poolSzie int
	if val, ok := o.Config["host"]; ok {
		host = val.(string)
	} else {
		host = "127.0.0.1"
	}
	if val, ok := o.Config["port"]; ok {
		port = val.(string)
	} else {
		port = "6379"
	}
	if val, ok := o.Config["password"]; ok {
		pass = val.(string)
	}
	if val, ok := o.Config["poolsize"]; ok {
		poolSzie = val.(int)
	} else {
		poolSzie = 10
	}
	c.handle = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
		PoolSize: poolSzie,
	})
	pong, err := c.handle.Ping().Result()
	if err != nil || pong != "PONG" {
		return fmt.Errorf("redis connect err: %s", err)
	}
	return nil
}

func init() {
	cache.Register("redis", &Redis{})
}
