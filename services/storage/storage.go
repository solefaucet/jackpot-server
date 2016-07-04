package storage

import "github.com/solefaucet/jackpot-server/models"

// Storage defines interface that one should implement for data persistence
type Storage interface {
	GetLatestBlock() (models.Block, error)

	SaveBlockAndTransactions(models.Block, []models.Transaction) error
}
