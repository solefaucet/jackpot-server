package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func upsertGame(tx *sqlx.Tx, hash string, height int64, totalAmount float64, blockCreatedAt time.Time) error {
	sql := "INSERT INTO `games` (`hash`, `height`, `total_amount`, `game_of`) VALUES (:hash, :height, :total_amount, :game_of) ON DUPLICATE KEY UPDATE `hash` = :hash, `height` = :height, `total_amount` = `total_amount` + :total_amount"
	_, err := tx.NamedExec(sql, map[string]interface{}{
		"hash":         hash,
		"total_amount": totalAmount,
		"height":       height,
		"game_of":      blockCreatedAt.Truncate(time.Hour),
	})
	if err != nil {
		return fmt.Errorf("upsert game error: %#v", err)
	}

	return nil
}
