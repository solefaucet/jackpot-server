package utils

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMust(t *testing.T) {
	Convey("Given nil error", t, func() {
		var err error
		var i int64 = 1

		Convey("When request must", func() {
			actual := Must(i, err)

			Convey("It should equal", func() {
				So(actual, ShouldEqual, 1)
			})
		})
	})

	Convey("Given non-nil error", t, func() {
		err := errors.New("")

		Convey("When request must", func() {
			f := func() { Must(nil, err) }

			Convey("It should panic", func() {
				So(f, ShouldPanic)
			})
		})
	})
}
