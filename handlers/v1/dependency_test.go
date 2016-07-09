package v1

import (
	"time"

	"github.com/solefaucet/jackpot-server/models"
)

func mockDependencyGetGamesWithin(games []models.Game, err error) dependencyGetGamesWithin {
	return func(_, _ time.Time) ([]models.Game, error) {
		return games, err
	}
}

func mockDependencyGetTransactionsWithin(transactions []models.Transaction, err error) dependencyGetTransactionsWithin {
	return func(_, _ time.Time) ([]models.Transaction, error) {
		return transactions, err
	}
}
