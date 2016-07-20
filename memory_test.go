package cache

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

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
			err := c.Set("test", "1", 2)
			So(err, ShouldBeNil)
		})

		Convey("get", func() {
			var v string
			c.Get("test", &v)
			So(v, ShouldEqual, "1")
		})

		Convey("get expried", func() {
			time.Sleep(time.Second * 2)
			var v string
			c.Get("test", &v)
			So(v, ShouldBeEmpty)
		})

		Convey("set 10000", func() {
			for i := 0; i < 10000; i++ {
				err = c.Set("test", "1", 10)
			}
			So(err, ShouldBeNil)
		})

		Convey("get 10000", func() {
			var v string
			for i := 0; i < 10000; i++ {
				v = ""
				c.Get("test", &v)
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
			So(v, ShouldEqual, 2)
			v, err = c.Incr("test1")
			So(v, ShouldEqual, 1)
			v, err = c.Decr("test2")
			So(err, ShouldNotBeNil)
			err = c.Delete("test")
			So(err, ShouldBeNil)
		})

		Convey("gc", func() {
			var v string
			for i := 0; i <= 11000; i++ {
				key := "test" + strconv.Itoa(i)
				err = c.Set(key, "01234567890123456789", 10)
			}
			So(err, ShouldBeNil)
			c.Get("test10000", &v)
			So(v, ShouldEqual, "01234567890123456789")
			v = ""
			c.Get("test6", &v)
			So(v, ShouldBeEmpty)
		})

		Convey("flush", func() {
			err := c.Flush()
			So(err, ShouldBeNil)
		})

		Convey("exists", func() {
			c.Incr("test1")
			ok := c.Exist("test1")
			So(ok, ShouldBeTrue)
			ok = c.Exist("testNotExist")
			So(ok, ShouldBeFalse)
		})

		Convey("large item", func() {
			v := strings.Repeat("A", 1024*1025)
			err := c.Set("test", v, 30)
			So(err, ShouldNotBeNil)
			v = strings.Repeat("A", 1024*513)
			err = c.Set("test2", v, 30)
			err = c.Set("test3", v, 30)
			So(err, ShouldBeNil)
			err = c.Delete("test3")
			So(err, ShouldBeNil)
		})
	})
}

func BenchmarkCacheMemorySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testCache.Set(fmt.Sprintf("test%d", i), 1, 1800)
	}
}

func BenchmarkCacheMemoryGet(b *testing.B) {
	var v string
	for i := 0; i < b.N; i++ {
		testCache.Get(fmt.Sprintf("test%d", i), &v)
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
