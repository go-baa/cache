# cache [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/cache) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/cache/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/cache.svg?style=flat-square)](https://travis-ci.org/go-baa/cache) [![Coverage Status](http://img.shields.io/coveralls/go-baa/cache.svg?style=flat-square)](https://coveralls.io/r/go-baa/cache)

baa module providers a cache management.

## Features

- multi storage support: memory, file, memcache, redis
- Get/Set/Incr/Decr/Delete/Exist/Flush/Start

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
        Prefix:   "MyApp",
		Adapter:  "memory",
		Config:   map[string]string{},
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

**Name**

``string``

the cache name

**Prefix**

``string``

the cache key prefix, used for isolate different cache instance/app.

**Adapter**

``string``

the cache adapter name, choose support adapter: memory, file, memcache, redis.

**Config**

``map[string]interface{}``

the cache adapter config, use a dict, values was diffrent with adapter.

### Adapter Memory

**bytesLimit**

``int64``

set the memory cache memory limit, default is 128m

**Usage**

```
app.SetDI("cache", cache.New(cache.Options{
    Name:     "cache",
    Adapter:  "memory",
    Config:   map[string]interface{}{
        "bytesLimit": int64(128 * 1024 * 1024), // 128m
    },
}))
```


