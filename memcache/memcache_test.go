package memcache

import (
	"encoding/gob"
	"fmt"
	"testing"
	"time"

	"github.com/go-baa/cache"
	. "github.com/smartystreets/goconvey/convey"
)

// init a global cacher
var c cache.Cacher

func TestCacheMemory1(t *testing.T) {
	Convey("cache redis", t, func() {
		Convey("set", func() {
			err := c.Set("test", "1", 2)
			So(err, ShouldBeNil)
		})

		Convey("get", func() {
			var v string
			c.Get("test", &v)
			So(v, ShouldEqual, "1")
		})

		Convey("get gc", func() {
			var v string
			time.Sleep(time.Second * 3)
			c.Get("test", &v)
			So(v, ShouldBeEmpty)
		})

		Convey("set struct", func() {
			type b struct {
				Name string
			}
			gob.Register(b{})
			v1 := b{"test"}
			err := c.Set("test", v1, 6)
			So(err, ShouldBeNil)
			var v2 b
			c.Get("test", &v2)
			So(v2.Name, ShouldEqual, v1.Name)
		})

		Convey("incr/decr", func() {
			c.Set("test", 1, 10)
			v, err := c.Incr("test")
			v, err = c.Incr("test")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 3)
			v, err = c.Decr("test")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 2)
		})
	})
}

func BenchmarkCacheMemorySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Set(fmt.Sprintf("test%d", i), 1, 1800)
	}
}

func BenchmarkCacheMemoryGet(b *testing.B) {
	var v string
	for i := 0; i < b.N; i++ {
		c.Get(fmt.Sprintf("test%d", i), v)
	}
}

func init() {
	c = cache.New(cache.Options{
		Name:    "test",
		Adapter: "memcache",
		Config: map[string]interface{}{
			"host": "127.0.0.1",
			"port": "11211",
		},
	})
}
