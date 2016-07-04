package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/jackpot-server/models"
)

func saveTransactions(tx *sqlx.Tx, transactions []models.Transaction) error {
	stmt, _ := tx.PrepareNamed("INSERT INTO `transactions` (`address`, `amount`, `tx_id`, `hash`, `block_created_at`) VALUES (:address, :amount, :tx_id, :hash, :block_created_at)")
	defer stmt.Close()

	for _, v := range transactions {
		if _, err := stmt.Exec(v); err != nil {
			return fmt.Errorf("save transaction error: %#v", err)
		}
	}

	return nil
}
