package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/jackpot-server/models"
)

func saveTransactions(tx *sqlx.Tx, transactions []models.Transaction) error {
	stmt, _ := tx.PrepareNamed("INSERT INTO `transactions` (`address`, `amount`, `tx_id`, `hash`, `block_created_at`, `game_of`, `confirmations`) VALUES (:address, :amount, :tx_id, :hash, :block_created_at, :game_of, :confirmations)")
	defer stmt.Close()

	for _, v := range transactions {
		if _, err := stmt.Exec(v); err != nil {
			return fmt.Errorf("save transactions error: %#v", err)
		}
	}

	return nil
}

// GetTransactionsWithin gets all transactions, filter by block_created_at = [start, end)
func (s Storage) GetTransactionsWithin(start, end time.Time) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := s.db.Select(&transactions, "SELECT * FROM `transactions` WHERE `block_created_at` >= ? AND `block_created_at` < ? ORDER BY `block_created_at` DESC", start, end)
	return transactions, err
}

// GetTransactionsByGameOf gets all transactions, filter by game_of
func (s Storage) GetTransactionsByGameOf(gameOf time.Time) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := s.db.Select(&transactions, "SELECT * FROM `transactions` WHERE `game_of` = ? ORDER BY `block_created_at` DESC", gameOf)
	return transactions, err
}

// UpdateTransactionConfirmationByID update confirmations by transaction id
func (s Storage) UpdateTransactionConfirmationByID(id int64, confirmations int64) error {
	sql := "UPDATE `transactions` SET `confirmations` = ? WHERE `id` = ?"
	_, err := s.db.Exec(sql, confirmations, id)
	return err
}
