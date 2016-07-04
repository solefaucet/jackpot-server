package mysql

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/jackpot-server/models"
)

func TestSaveTransactions(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save duplicated transactions", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				return saveTransactions(tx, []models.Transaction{
					{
						Address:        "addr",
						Amount:         10,
						TransactionID:  "id",
						Hash:           "hash",
						BlockCreatedAt: time.Now(),
					},
					{
						Address:        "addr",
						Amount:         10,
						TransactionID:  "id",
						Hash:           "hash",
						BlockCreatedAt: time.Now(),
					},
				})
			})

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save transactions", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				return saveTransactions(tx, []models.Transaction{
					{
						Address:        "addr",
						Amount:         10,
						TransactionID:  "id",
						Hash:           "hash",
						BlockCreatedAt: time.Now(),
					},
				})
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
