package configurer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONConfigure(t *testing.T) {
	Convey("Given a json configurer and some helpers", t, func() {
		c := &JSON{
			Source: []byte(`{"Foo":1.5,"Goo":true,"Hoo":"horton"}`),
		}
		var f struct{ Foo float64 }
		var g struct{ Goo bool }
		var h struct{ Hoo string }

		Convey("When configuring f", func() {
			err := c.Configure(&f)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The value should be set", func() { So(f.Foo, ShouldEqual, 1.5) })
		})

		Convey("When configuring g", func() {
			err := c.Configure(&g)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The value should be set", func() { So(g.Goo, ShouldEqual, true) })
		})

		Convey("When configuring h", func() {
			err := c.Configure(&h)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The value should be set", func() { So(h.Hoo, ShouldEqual, "horton") })
		})

		Convey("When the source is malformed", func() {
			c.Source = []byte(`{Foo:"}`)
			err := c.Configure(&f)
			Convey("There should be no error", func() { So(err, ShouldNotBeNil) })
		})
	})
}
