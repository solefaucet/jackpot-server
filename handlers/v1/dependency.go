package v1

import (
	"time"

	"github.com/solefaucet/jackpot-server/models"
)

type (
	dependencyGetGamesWithin        func(start, end time.Time) ([]models.Game, error)
	dependencyGetTransactionsWithin func(start, end time.Time) ([]models.Transaction, error)
)
