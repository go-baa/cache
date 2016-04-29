package cache

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// init a global cacher
var testCache Cacher
var err error

func TestCacheMemory1(t *testing.T) {
	Convey("cache memory", t, func() {
		c := New(Options{
			Name:    "test2",
			Adapter: "memory",
			Config: map[string]interface{}{
				"bytesLimit": int64(1024 * 1024), // 1MB
			},
		})

		Convey("set", func() {
			err := c.Set("test", "1", 6)
			So(err, ShouldBeNil)
		})

		Convey("get", func() {
			v := c.Get("test")
			So(v, ShouldEqual, "1")
		})

		Convey("set 10000", func() {
			for i := 0; i < 10000; i++ {
				err = c.Set("test", "1", 10)
			}
			So(err, ShouldBeNil)
		})

		Convey("get 10000", func() {
			var v interface{}
			for i := 0; i < 10000; i++ {
				v = c.Get("test")
			}
			So(v, ShouldEqual, "1")
		})

		Convey("set struct", func() {
			type b struct {
				Name string
			}
			gob.Register(b{})
			v1 := b{"test"}
			err := c.Set("test", v1, 6)
			So(err, ShouldBeNil)
			v2 := c.Get("test")
			So(v2.(b).Name, ShouldEqual, v1.Name)
		})

		Convey("incr/decr", func() {
			c.Set("test", 1, 10)
			v, err := c.Incr("test")
			v, err = c.Incr("test")
			So(err, ShouldBeNil)
			So(v.(int), ShouldEqual, 3)
			v, err = c.Decr("test")
			So(err, ShouldBeNil)
			So(v.(int), ShouldEqual, 2)
		})

		Convey("gc", func() {
			for i := 0; i <= 200000; i++ {
				key := "test" + strconv.Itoa(i)
				err = c.Set(key, "01234567890123456789", 6)
			}
			So(err, ShouldBeNil)
			v := c.Get("test200000")
			So(v, ShouldEqual, "01234567890123456789")
			v = c.Get("test100")
			So(v, ShouldBeNil)
		})
	})
}

func BenchmarkCacheMemorySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testCache.Set(fmt.Sprintf("test%d", i), 1, 1800)
	}
}

func BenchmarkCacheMemoryGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testCache.Get(fmt.Sprintf("test%d", i))
	}
}

func init() {
	testCache = New(Options{
		Name:    "test",
		Adapter: "memory",
		Config: map[string]interface{}{
			"bytesLimit": int64(1024 * 1024), // 1MB
		},
	})
}
