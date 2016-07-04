package wallet

import "time"

// Wallet defines interface that one should implement for blockchain manipulation
type Wallet interface {
	GetBlock(bestBlock bool, height int64) (*Block, error)
	GetReceivedSince(hash string, minConfirms int) ([]Transaction, error)
}

// Block _
type Block struct {
	Height         int64
	PrevHash       string
	Hash           string
	BlockCreatedAt time.Time
}

// Transaction _
type Transaction struct {
	Address        string
	Amount         float64
	TransactionID  string
	Hash           string
	BlockCreatedAt time.Time
}
