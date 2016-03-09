package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache1(t *testing.T) {
	Convey("cache", t, func() {
		c := New(Options{
			Name:     "test",
			Adapter:  "memory",
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
	})
}
