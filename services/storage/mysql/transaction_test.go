package mysql

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/jackpot-server/models"
)

func TestSaveTransactions(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save duplicated transactions", func() {
			err := s.SaveTransactions([]models.Transaction{
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

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save transactions", func() {
			err := s.SaveTransactions([]models.Transaction{
				{
					Address:        "addr",
					Amount:         10,
					TransactionID:  "id",
					Hash:           "hash",
					BlockCreatedAt: time.Now(),
				},
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetLatestTransactionBlockHash(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get latest block", func() {
			hash, err := s.GetLatestTransactionBlockHash()

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Hash should be empty", func() {
				So(hash, ShouldBeEmpty)
			})
		})
	})

	Convey("Given mysql storage with transaction data", t, func() {
		s := prepareDatabaseForTesting()
		s.SaveTransactions([]models.Transaction{
			{
				Address:        "addr",
				Amount:         10,
				TransactionID:  "id1",
				Hash:           "hash1",
				BlockCreatedAt: time.Now(),
			},
			{
				Address:        "addr",
				Amount:         10,
				TransactionID:  "id2",
				Hash:           "hash2",
				BlockCreatedAt: time.Now().Add(time.Hour),
			},
		})

		Convey("When get latest transaction block hash", func() {
			hash, _ := s.GetLatestTransactionBlockHash()

			Convey("Hash should be hash2", func() {
				So(hash, ShouldEqual, "hash2")
			})
		})
	})

	withClosedConn(t, "When get latest transaction block hash", func(s Storage) error {
		_, err := s.GetLatestTransactionBlockHash()
		return err
	})
}
