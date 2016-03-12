# cache
baa module providers a cache management.

## Features

- multi storage support: memory, file, memcache, redis, couchbase
- Get/Set/Delete/Exist/Flush/Start

## Getting Started

```
package main

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/cache"
)

func main() {
	// new app
	app := baa.New()

	// register cache
	app.SetDI("cache", cache.New(cache.Options{
		Name:     "cache",
		Adapter:  "memory",
		Config:   map[string]string{},
		Interval: 60,
	}))

	// router
	app.Get("/", func(c *baa.Context) {
		ca := c.DI("cache").(cache.Cacher)
		ca.Set("test", "baa", 10)
		v := ca.Get("test").(string)
		c.String(200, v)
	})

	// run app
	app.Run(":1323")
}
```

you should import cache adpater before use it, like this:

```
import(
    "github.com/go-baa/baa"
    "github.com/go-baa/cache"
    _ "github.com/go-baa/cache/memcache"
)
```

adapter ``memory`` has build in, do not need import.

## Configuration

### Common

** Name string **

the cache name

** Adapter string **

the cache adapter name, choose support adapter: memory, file, memcache, redis.

** Config map[string]string **

the cache adapter config, use a dict, values was diffrent with adapter.

** Interval int64 **

the cache gc interval time, second.

### Adapter Memory

** bytesLimit int64 **

set the memory cache memory limit, default is 128m

** usage **

```
app.SetDI("cache", cache.New(cache.Options{
    Name:     "cache",
    Adapter:  "memory",
    Config:   map[string]string{
        "bytesLimit": "134217728", // 128m
    },
    Interval: 60,
}))
```


