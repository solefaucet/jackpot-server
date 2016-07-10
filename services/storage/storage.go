package storage

import (
	"time"

	"github.com/solefaucet/jackpot-server/models"
)

// Storage defines interface that one should implement for data persistence
type Storage interface {
	// block
	GetLatestBlock() (models.Block, error)

	// transaction
	GetTransactionsWithin(start, end time.Time) ([]models.Transaction, error)
	GetUnconfirmedTransactions(confirmations int64) ([]models.Transaction, error)
	GetTransactionsByGameOf(gameOf time.Time) ([]models.Transaction, error)
	UpdateTransactionConfirmationByID(id int64, confirmations int64) error

	// game
	GetGamesWithin(start, end time.Time) ([]models.Game, error)
	GetDrawingNeededGames() ([]models.Game, error)
	UpdateGameToEndedStatus(models.Game) error

	// batch
	SaveBlockAndTransactions(time.Time, models.Block, []models.Transaction, *models.Game) error
}
