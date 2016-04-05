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
	handle *redis.Client
}

// Exist check key is exist
func (c *Redis) Exist(key string) bool {
	ok, err := c.handle.Exists(key).Result()
	if err == nil && ok {
		return true
	}
	return false
}

// Get returns value for given key
func (c *Redis) Get(key string) interface{} {
	v, err := c.handle.Get(key).Bytes()
	if err != nil {
		return nil
	}
	item, err := cache.ItemBinary(v).Item()
	if err != nil || item == nil {
		return nil
	}
	return item.Val
}

// Set set value for given key
func (c *Redis) Set(key string, v interface{}, ttl int64) error {
	item := cache.NewItem(v, ttl)
	b, err := item.Bytes()
	if err != nil {
		return err
	}
	_, err = c.handle.Set(key, []byte(b), time.Second*time.Duration(ttl)).Result()
	return err
}

// Delete delete the key
func (c *Redis) Delete(key string) error {
	return c.handle.Del(key).Err()
}

// Flush flush cacher
func (c *Redis) Flush() error {
	return c.handle.FlushDb().Err()
}

// Start new a cacher and start service
func (c *Redis) Start(o cache.Options) error {
	c.Name = o.Name
	var host, port, pass string
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
	c.handle = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
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
