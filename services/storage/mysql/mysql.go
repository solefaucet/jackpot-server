package mysql

import (
	_ "github.com/go-sql-driver/mysql" // is needed for mysql driver registeration
	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/jackpot-server/models"
	"github.com/solefaucet/jackpot-server/services/storage"
)

// Storage implements Storage interface for data storage
type Storage struct {
	db *sqlx.DB
}

var _ storage.Storage = Storage{}

// New returns a Storage with data source name
func New(dsn string) Storage {
	return Storage{
		db: sqlx.MustConnect("mysql", dsn).Unsafe(),
	}
}

// SetMaxOpenConns alias sql.DB.SetMaxOpenConns
func (s *Storage) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}

// SetMaxIdleConns alias sql.DB.SetMaxIdleConns
func (s *Storage) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

func (s Storage) withTx(f func(*sqlx.Tx) error) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// SaveBlockAndTransactions save block and transactions
func (s Storage) SaveBlockAndTransactions(block models.Block, transactions []models.Transaction) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		if err := saveBlock(tx, block); err != nil {
			return err
		}

		if err := saveTransactions(tx, transactions); err != nil {
			return err
		}

		totalAmount := 0.0
		for _, v := range transactions {
			totalAmount += v.Amount
		}
		if err := upsertGame(tx, block.Hash, block.Height, totalAmount, block.BlockCreatedAt); err != nil {
			return err
		}

		return nil
	})
}
