package mysql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/jackpot-server/models"
)

// GetLatestTransactionBlockHash gets latest received transaction block hash
func (s Storage) GetLatestTransactionBlockHash() (string, error) {
	var hash string
	err := s.db.Get(&hash, "SELECT `hash` FROM `transactions` ORDER BY `block_created_at` DESC LIMIT 1")

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		err = fmt.Errorf("get latest transaction block hash error: %#v", err)
	}

	return hash, err
}

// SaveTransactions saves transactions to db
func (s Storage) SaveTransactions(transactions []models.Transaction) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		stmt, _ := tx.PrepareNamed("INSERT INTO `transactions` (`address`, `amount`, `tx_id`, `hash`, `block_created_at`) VALUES (:address, :amount, :tx_id, :hash, :block_created_at)")
		defer stmt.Close()

		for _, v := range transactions {
			if _, err := stmt.Exec(v); err != nil {
				return fmt.Errorf("save transaction error: %#v", err)
			}
		}

		return nil
	})
}
