package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
	w "github.com/solefaucet/jackpot-server/services/wallet"
)

var (
	blockHeightChan        = make(chan int64, 2)
	previousBlockCreatedAt = time.Time{}
)

func initWork() {
	go fetchAndSave()

	// get latest block from db
	block, err := storage.GetLatestBlock()

	if err == jerrors.ErrNotFound {
		blockHeightChan <- -1
		return
	}

	if err != nil {
		logger.Panicf("fail to get latest block from database: %#v\n", err)
		return
	}

	previousBlockCreatedAt = block.BlockCreatedAt
	blockHeightChan <- block.Height + 1
}

func fetchAndSave() {
	for {
		height := <-blockHeightChan
		saveBlockAndTransactions(height)
	}
}

func saveBlockAndTransactions(height int64) {
	var err error
	defer func() {
		if err != nil {
			time.Sleep(time.Minute)
		}

		blockHeightChan <- height
	}()

	entry := logrus.WithFields(logrus.Fields{
		"event":        models.LogEventSaveBlockAndTransactions,
		"block_height": height,
	})

	// get new block from blockchain
	bestBlock := height < 0
	block, err := wallet.GetBlock(bestBlock, height)
	if err == jerrors.ErrNoNewBlock {
		entry.Info("no new block ahead")
		return
	}

	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to get block from blockchain")
		return
	}

	gameOf := block.BlockCreatedAt.Truncate(config.Jackpot.Duration)
	entry.WithFields(logrus.Fields{
		"previous_hash": block.PrevHash,
		"height":        block.Height,
		"game_of":       gameOf,
	})

	// get receive transactions
	transactions, err := wallet.GetReceivedSince(block.PrevHash, block.Hash)
	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to list transactions from blockchain")
		return
	}

	// check if it's time to find out the winner
	previousGameOf := previousBlockCreatedAt.Truncate(config.Jackpot.Duration)
	var updatedGame *models.Game
	if previousGameOf.Add(config.Jackpot.Duration).Equal(gameOf) {
		updatedGame = &models.Game{
			Hash:   block.Hash,
			Height: block.Height,
			GameOf: previousGameOf,
		}
	}

	// save block, transactions
	if err := storage.SaveBlockAndTransactions(
		gameOf,
		walletBlockToModelBlock(block),
		walletTxsToModelTxs(gameOf, transactions),
		updatedGame,
	); err != nil {
		entry.WithField("error", err.Error()).Error("fail to save block and transactions to db")
		return
	}

	entry.Info("save block and transactions successfully")
	height = block.Height + 1
	previousBlockCreatedAt = block.BlockCreatedAt
}

func walletTxsToModelTxs(gameOf time.Time, txs []w.Transaction) []models.Transaction {
	transactions := make([]models.Transaction, len(txs))
	for i, v := range txs {
		transactions[i] = models.Transaction{
			Address:        v.Address,
			Amount:         v.Amount,
			TransactionID:  v.TransactionID,
			Hash:           v.Hash,
			Confirmations:  v.Confirmations,
			GameOf:         gameOf,
			BlockCreatedAt: v.BlockCreatedAt,
		}
	}
	return transactions
}

func walletBlockToModelBlock(blockchainBlock *w.Block) models.Block {
	return models.Block{
		Hash:           blockchainBlock.Hash,
		Height:         blockchainBlock.Height,
		BlockCreatedAt: blockchainBlock.BlockCreatedAt,
	}
}

	}
}
