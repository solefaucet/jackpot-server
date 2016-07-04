package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/jackpot-server/models"
)

// test helpers
func execCommand(cmd string) {
	c := exec.Command("sh", "-c", "-i", cmd)
	if err := c.Run(); err != nil {
		log.Fatalf("execute command %v, error: %v", cmd, err)
	}
}

func prepareDatabaseForTesting() Storage {
	execCommand(`mysql -uroot -e 'drop database if exists jackpot_test;'`)
	execCommand(`mysql -uroot -e 'create database jackpot_test character set utf8;'`)
	dir, _ := os.Getwd()
	paths := strings.Split(dir, string(os.PathSeparator))
	projBasePath := strings.Join(paths[:len(paths)-3], string(os.PathSeparator)) // cannot come up with any better way
	execCommand(fmt.Sprintf(`cd %s && goose -env test up`, projBasePath))

	dsn := "root:@/jackpot_test?parseTime=true"
	s := New(dsn)
	s.SetMaxOpenConns(4)
	s.SetMaxIdleConns(4)
	return s
}

func resetDatabase() {
	execCommand(`mysql -uroot -e 'drop database if exists jackpot_test;'`)
}

func withClosedConn(t *testing.T, description string, f func(Storage) error) {
	Convey("Given mysql storage with closed connection", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Close()

		Convey(description, func() {
			err := f(s)

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestWithTx(t *testing.T) {
	withClosedConn(t, "When begin transaction", func(s Storage) error {
		return s.withTx(func(*sqlx.Tx) error {
			return nil
		})
	})

	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		expected := fmt.Errorf("errored")

		Convey("When execute errored sql", func() {
			actual := s.withTx(func(*sqlx.Tx) error {
				return expected
			})

			Convey("It should resemble", func() {
				So(actual, ShouldResemble, expected)
			})
		})
	})
}

func TestSaveBlockAndTransactions(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When save block and transactions", func() {
			err := s.SaveBlockAndTransactions(
				models.Block{Hash: "hash", Height: 1, BlockCreatedAt: time.Now()},
				[]models.Transaction{
					{
						Address:        "addr",
						Amount:         10,
						TransactionID:  "id",
						Hash:           "hash",
						BlockCreatedAt: time.Now(),
					},
				},
			)

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
