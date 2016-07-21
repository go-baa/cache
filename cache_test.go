package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheGob1(t *testing.T) {
	Convey("cache Gob", t, func() {
		item := NewItem("1", 6)
		b, err := item.Encode()
		item2, err2 := b.Item()

		Convey("encode", func() {
			So(err, ShouldBeNil)
		})

		Convey("decode", func() {
			So(err2, ShouldBeNil)
			So(item2.Val, ShouldEqual, "1")
		})
	})
}

func TestCache1(t *testing.T) {
	Convey("cache", t, func() {
		c := New(Options{
			Name:    "test",
			Adapter: "memory",
		})

		Convey("set", func() {
			err := c.Set("test", "1", 6)
			So(err, ShouldBeNil)
		})

		Convey("get", func() {
			var v string
			c.Get("test", &v)
			So(v, ShouldEqual, "1")
		})
	})
}
