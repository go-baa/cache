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
        ca := c.DI("cache")
        ca.Set("test", "baa")
        v := ca.Get("test")
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

memory adapter has build in, do not need import.
