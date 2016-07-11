package v1

import (
	"time"

	"github.com/solefaucet/jackpot-server/models"
)

func mockDependencyGetGames(games []models.Game, err error) dependencyGetGames {
	return func(_, _ int64) ([]models.Game, error) {
		return games, err
	}
}

func mockDependencyGetTransactionsByGameOfs(transactions []models.Transaction, err error) dependencyGetTransactionsByGameOfs {
	return func(...time.Time) ([]models.Transaction, error) {
		return transactions, err
	}
}
