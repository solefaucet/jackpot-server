package v1

import (
	"time"

	"github.com/solefaucet/jackpot-server/models"
)

type (
	dependencyGetGames                 func(limit, offset int64) ([]models.Game, error)
	dependencyGetTransactionsByGameOfs func(gameOfs ...time.Time) ([]models.Transaction, error)
)
