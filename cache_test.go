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
	c := New(Options{
		Name:    "test",
		Adapter: "memory",
	})

	Convey("cache", t, func() {

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

func TestCacheMulti(t *testing.T) {
	Convey("test multiple adapter", t, func() {
		c1 := New(Options{
			Name:    "test1",
			Adapter: "memory",
			Config: map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     "6379",
				"password": "",
				"poolsize": 10,
			},
		})
		c2 := New(Options{
			Name:    "test2",
			Adapter: "memory",
			Config: map[string]interface{}{
				"host":     "10.1.1.31",
				"port":     "6379",
				"password": "",
				"poolsize": 10,
			},
		})
		So(c1, ShouldNotEqual, c2)
	})
}
