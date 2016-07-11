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

// GetUnconfirmedTransactions gets all unconfirmed transactions
func (s Storage) GetUnconfirmedTransactions(confirmations int64) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	err := s.db.Select(&transactions, "SELECT * FROM `transactions` WHERE `confirmations` < ? ORDER BY `block_created_at` DESC", confirmations)
	return transactions, err
}

// GetTransactionsByGameOfs gets all transactions, filter by game_of
func (s Storage) GetTransactionsByGameOfs(gameOfs ...time.Time) ([]models.Transaction, error) {
	sql, args, err := sqlx.In(
		"SELECT * FROM `transactions` WHERE `game_of` IN (?) ORDER BY `block_created_at` DESC",
		gameOfs,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to build sql with in: %v", err)
	}

	transactions := []models.Transaction{}
	err = s.db.Select(&transactions, sql, args...)
	return transactions, err
}

// UpdateTransactionConfirmationByID update confirmations by transaction id
func (s Storage) UpdateTransactionConfirmationByID(id int64, confirmations int64) error {
	sql := "UPDATE `transactions` SET `confirmations` = ? WHERE `id` = ?"
	_, err := s.db.Exec(sql, confirmations, id)
	return err
}
