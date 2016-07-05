package mysql

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
)

func TestSaveBlock(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save block", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				return saveBlock(tx, models.Block{Hash: "hash", Height: 1, BlockCreatedAt: time.Now()})
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save block with commited connection", func() {
			err := s.withTx(func(tx *sqlx.Tx) error {
				tx.Commit()
				return saveBlock(tx, models.Block{Hash: "hash", Height: 1, BlockCreatedAt: time.Now()})
			})

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetLatestBlock(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get latest block", func() {
			_, err := s.GetLatestBlock()

			Convey("Error should be not found", func() {
				So(err, ShouldEqual, jerrors.ErrNotFound)
			})
		})
	})

	Convey("Given mysql storage with block data", t, func() {
		s := prepareDatabaseForTesting()
		s.withTx(func(tx *sqlx.Tx) error {
			return saveBlock(tx, models.Block{Hash: "hash", Height: 1, BlockCreatedAt: time.Now()})
		})
		s.withTx(func(tx *sqlx.Tx) error {
			return saveBlock(tx, models.Block{Hash: "hash2", Height: 2, BlockCreatedAt: time.Now()})
		})

		Convey("When get latest block", func() {
			block, _ := s.GetLatestBlock()

			Convey("Hash should be hash2", func() {
				So(block.Hash, ShouldEqual, "hash2")
			})
		})
	})

	withClosedConn(t, "When get latest block", func(s Storage) error {
		_, err := s.GetLatestBlock()
		return err
	})
}
