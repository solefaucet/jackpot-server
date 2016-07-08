package wallet

import "time"

// Wallet defines interface that one should implement for blockchain manipulation
type Wallet interface {
	GetBlock(bestBlock bool, height int64) (*Block, error)
	GetReceivedSince(prevHash, curHash string) ([]Transaction, error)
	SendToAddress(address string, amount float64) (string, error)
	GetDestAddress() (string, error)
	GetConfirmationsFromTxID(txID string) (int64, error)
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
	Confirmations  int64
	BlockCreatedAt time.Time
}
