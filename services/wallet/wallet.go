package wallet

import "time"

// Wallet defines interface that one should implement for blockchain manipulation
type Wallet interface {
	GetBlock(bestBlock bool, height int64) (*Block, error)
}

// Block _
type Block struct {
	Height         int64
	Hash           string
	BlockCreatedAt time.Time
}
