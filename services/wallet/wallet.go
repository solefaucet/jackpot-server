package wallet

import "time"

// Wallet defines interface that one should implement for blockchain manipulation
type Wallet interface {
	GetBestBlockHash() (string, error)
	GetBlockHash(height int64) (string, error)
	GetBlock(hash string) (Block, error)
}

// Block _
type Block struct {
	Height         int64
	Hash           string
	BlockCreatedAt time.Time
}
