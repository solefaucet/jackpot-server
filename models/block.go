package models

import "time"

// Block model
type Block struct {
	ID             int64     `db:"id"`
	Hash           string    `db:"hash"`
	Height         int32     `db:"height"`
	BlockCreatedAt time.Time `db:"block_created_at"`
	CreatedAt      time.Time `db:"created_at"`
}
