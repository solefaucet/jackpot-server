package storage

import "github.com/solefaucet/jackpot-server/models"

// Storage defines interface that one should implement for data persistence
type Storage interface {
	// Blocks
	GetLatestBlock() (models.Block, error)
	SaveBlock(models.Block) error
}
