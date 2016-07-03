package models

import "time"

// Transaction model
type Transaction struct {
	ID             int64     `db:"id"`
	Address        string    `db:"address"`
	Amount         float64   `db:"amount"`
	TransactionID  string    `db:"tx_id"`
	Hash           string    `db:"hash"`
	BlockCreatedAt time.Time `db:"block_created_at"`
	CreatedAt      time.Time `db:"created_at"`
}
