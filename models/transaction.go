package models

import "time"

// Transaction model
type Transaction struct {
	ID             int64     `db:"id"`
	Address        string    `db:"address"`
	Amount         float64   `db:"amount"`
	TransactionID  string    `db:"tx_id"`
	Confirmations  int64     `db:"confirmations"`
	Hash           string    `db:"hash"`
	GameOf         time.Time `db:"game_of"`
	BlockCreatedAt time.Time `db:"block_created_at"`
	CreatedAt      time.Time `db:"created_at"`
}
