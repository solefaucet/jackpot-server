package mysql

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUpsertGame(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When upsert game", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				return upsertGame(tx, "hash", 1, 0.2, time.Now().UTC())
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When upsert game with commited connection", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				tx.Commit()
				return upsertGame(tx, "hash", 1, 0.2, time.Now().UTC())
			})

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
