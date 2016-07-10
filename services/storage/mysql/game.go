package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/jackpot-server/models"
)

func upsertGame(tx *sqlx.Tx, gameOf time.Time, hash string, height int64, totalAmount float64) error {
	sql := "INSERT INTO `games` (`hash`, `height`, `total_amount`, `game_of`) VALUES (:hash, :height, :total_amount, :game_of) ON DUPLICATE KEY UPDATE `hash` = :hash, `height` = :height, `total_amount` = `total_amount` + :total_amount"
	_, err := tx.NamedExec(sql, map[string]interface{}{
		"hash":         hash,
		"total_amount": totalAmount,
		"height":       height,
		"game_of":      gameOf,
	})
	if err != nil {
		return fmt.Errorf("upsert game error: %#v", err)
	}

	return nil
}

func updateGameToDrawingNeededStatus(tx *sqlx.Tx, game *models.Game) error {
	if game == nil {
		return nil
	}

	sql := "UPDATE `games` SET `hash` = ?, `height` = ?, `status` = ? WHERE `game_of` = ? AND `status` = ?"
	result, err := tx.Exec(sql, game.Hash, game.Height, models.GameStatusDrawingNeeded, game.GameOf, models.GameStatusPending)
	if err != nil {
		return fmt.Errorf("update game to drawing needed status error: %#v", err)
	}

	if affect, _ := result.RowsAffected(); affect != 1 {
		return fmt.Errorf("update game to drawing needed status affected row not 1 but %v", affect)
	}

	return nil
}

// GetGamesWithin gets all games, filter by game_of = [start, end)
func (s Storage) GetGamesWithin(start, end time.Time) ([]models.Game, error) {
	games := []models.Game{}
	err := s.db.Select(&games, "SELECT * FROM `games` WHERE `game_of` >= ? AND `game_of` < ? ORDER BY `game_of` DESC", start, end)
	return games, err
}

// GetDrawingNeededGames gets all drawing games
func (s Storage) GetDrawingNeededGames() ([]models.Game, error) {
	games := []models.Game{}
	err := s.db.Select(&games, "SELECT * FROM `games` WHERE `status` = ? ORDER BY `game_of` ASC", models.GameStatusDrawingNeeded)
	return games, err
}

// UpdateGameToEndedStatus updates game status to ended
func (s Storage) UpdateGameToEndedStatus(game models.Game) error {
	sql := "UPDATE `games` SET `address` = ?, `win_amount` = ?, `fee` = ?, `tx_id` = ?, `status` = ? WHERE `game_of` = ? AND `status` = ?"
	result, err := s.db.Exec(sql, game.Address, game.WinAmount, game.Fee, game.TransactionID, models.GameStatusEnded, game.GameOf, models.GameStatusDrawingNeeded)
	if err != nil {
		return fmt.Errorf("update game to ended status error: %#v", err)
	}

	if affect, _ := result.RowsAffected(); affect != 1 {
		return fmt.Errorf("update game to ended status affected row not 1 but %v", affect)
	}

	return nil
}
