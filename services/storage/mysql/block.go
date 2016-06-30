package mysql

import (
	"database/sql"
	"fmt"

	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
)

// GetLatestBlock gets latest models.Block
func (s Storage) GetLatestBlock() (models.Block, error) {
	block := models.Block{}
	err := s.db.Get(&block, "SELECT * FROM `blocks` ORDER BY `height` DESC LIMIT 1")

	if err != nil {
		if err == sql.ErrNoRows {
			return block, jerrors.ErrNotFound
		}

		return block, fmt.Errorf("get latest block error: %#v", err)
	}

	return block, nil
}

// SaveBlock saves models.Block
func (s Storage) SaveBlock(block models.Block) error {
	_, err := s.db.NamedExec("INSERT INTO `blocks` (`hash`, `height`, `block_created_at`) VALUES (:hash, :height, :block_created_at)", block)
	if err != nil {
		return fmt.Errorf("save block error: %#v", err)
	}

	return nil
}
