package cache

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheGob1(t *testing.T) {
	Convey("cache Gob", t, func() {
		item := &Item{
			Val:     "1",
			Created: time.Now().Unix(),
			TTL:     6,
		}
		item2 := new(Item)
		b, err := EncodeGob(item)
		err2 := DecodeGob(b, item2)

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
