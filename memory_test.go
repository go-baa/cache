package cache

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheMemory1(t *testing.T) {
	Convey("cache memory", t, func() {
		c := New(Options{
			Name:    "test",
			Adapter: "memory",
			Config: map[string]interface{}{
				"bytesLimit": int64(1024), // 1KB
			},
			Interval: 60,
		})

		Convey("set", func() {
			err := c.Set("test", "1", 6)
			So(err, ShouldBeNil)
		})

		Convey("get", func() {
			v := c.Get("test")
			So(v, ShouldEqual, "1")
		})

		Convey("gc", func() {
			for i := 0; i <= 100; i++ {
				key := "test" + strconv.Itoa(i)
				err := c.Set(key, i, 6)
				So(err, ShouldBeNil)
			}
			v := c.Get("test100")
			So(v, ShouldEqual, 100)
			v = c.Get("test1")
			So(v, ShouldBeNil)
		})
	})
}
