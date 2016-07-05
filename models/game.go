package models

import "time"

// Game model
type Game struct {
	ID            int64     `db:"id"`
	Hash          string    `db:"hash"`
	Height        int64     `db:"height"`
	Address       string    `db:"address"`
	WinAmount     float64   `db:"win_amount"`
	TotalAmount   float64   `db:"total_amount"`
	Fee           float64   `db:"fee"`
	TransactionID string    `db:"tx_id"`
	GameOf        time.Time `db:"game_of"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
