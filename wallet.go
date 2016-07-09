package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
	w "github.com/solefaucet/jackpot-server/services/wallet"
	"github.com/solefaucet/jackpot-server/utils"
)

var (
	blockHeightChan        = make(chan int64, 2)
	previousBlockCreatedAt = time.Time{}
)

func initWork() {
	go fetchAndSave()
	go processAndUpdate()
	go updateConfirmationsJob()

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

func processAndUpdate() {
func updateConfirmationsJob() {
	for {
		updateConfirmations()
		time.Sleep(time.Minute)
	}
}

	for {
		processGames()
		time.Sleep(time.Minute)
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

func updateConfirmations() {
	minConfirmations := config.Wallet.MinConfirms

	entry := logrus.WithFields(logrus.Fields{
		"event":             models.LogEventUpdateConfirmations,
		"min_confirmations": minConfirmations,
	})

	transactions, err := storage.GetUnconfirmedTransactions(minConfirmations)
	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to get unconfirmed transactions")
		return
	}

	for _, transaction := range transactions {
		confirmations, err := wallet.GetConfirmationsFromTxID(transaction.TransactionID)
		if err != nil {
			entry.WithFields(logrus.Fields{
				"error": err.Error(),
				"tx_id": transaction.TransactionID,
			}).Error("fail to get confirmations by tx id")
			return
		}

		if err := storage.UpdateTransactionConfirmationByID(transaction.ID, confirmations); err != nil {
			entry.WithFields(logrus.Fields{
				"error": err.Error(),
				"tx_id": transaction.TransactionID,
			}).Error("fail to update transaction confirmation")
			return
		}
	}
}

func processGames() {
	entry := logrus.WithFields(logrus.Fields{
		"event": models.LogEventProcessGames,
	})

	games, err := storage.GetProcessingGames()
	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to get processing games")
		return
	}

	for _, game := range games {
		transactions, err := storage.GetTransactionsByGameOf(game.GameOf)
		if err != nil {
			entry.WithFields(logrus.Fields{
				"error":   err.Error(),
				"game_of": game.GameOf,
			}).Error("fail to get transactions by game_of")
			return
		}

		processingNeeded := true
		for _, transaction := range transactions {
			confirmations, err := wallet.GetConfirmationsFromTxID(transaction.TransactionID)
			if err != nil {
				entry.WithFields(logrus.Fields{
					"error": err.Error(),
					"tx_id": transaction.TransactionID,
				}).Error("fail to get confirmations by tx id")
				return
			}

			if config.Wallet.MinConfirms > confirmations {
				processingNeeded = false
			}

			if err := storage.UpdateTransactionConfirmationByID(transaction.ID, confirmations); err != nil {
				entry.WithFields(logrus.Fields{
					"error": err.Error(),
					"tx_id": transaction.TransactionID,
				}).Error("fail to update transaction confirmation")
				return
			}
		}

		if processingNeeded {
			winnerAddress, transactionID, winAmount, fee, err := findWinnerAndSendCoins(game.GameOf, game.Hash)
			if err != nil {
				entry.WithFields(logrus.Fields{
					"game_of": game.GameOf,
					"hash":    game.Hash,
					"error":   err.Error(),
				}).Error("fail to find winner and send coins")
				return
			}

			g := models.Game{
				Address:       winnerAddress,
				WinAmount:     winAmount,
				Fee:           fee,
				TransactionID: transactionID,
				GameOf:        game.GameOf,
			}
			if err := storage.UpdateGameToEndedStatus(g); err != nil {
				entry.WithFields(logrus.Fields{
					"winner_address": winnerAddress,
					"win_amount":     winAmount,
					"fee":            fee,
					"tx_id":          transactionID,
					"game_of":        game.GameOf,
				}).Panic("fail to update game status to ended")
				return
			}
		}
	}
}

func findWinnerAndSendCoins(gameOf time.Time, hash string) (winnerAddress, transactionID string, winAmount, fee float64, err error) {
	transactions, err := storage.GetTransactionsByGameOf(gameOf)
	if err != nil {
		return
	}

	// no transactions, no winner
	if len(transactions) == 0 {
		return
	}

	totalAmount := totalAmountOfTransactions(transactions)
	fee = totalAmount * config.Jackpot.TransactionFee
	winAmount = totalAmount - fee

	winnerAddress, err = utils.FindWinner(transactions, hash)
	if err != nil {
		return
	}

	transactionID, err = wallet.SendToAddress(winnerAddress, winAmount)
	if err != nil {
		return
	}

	return
}

func totalAmountOfTransactions(transactions []models.Transaction) float64 {
	totalAmount := 0.0
	for _, tx := range transactions {
		totalAmount += tx.Amount
	}
	return totalAmount
}
