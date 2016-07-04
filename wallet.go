package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/models"
	w "github.com/solefaucet/jackpot-server/services/wallet"
)

var blockHeightChan = make(chan int64, 2)

func initWork() {
	go work()

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

	blockHeightChan <- block.Height + 1
}

func work() {
	for {
		select {
		case height := <-blockHeightChan:
			saveBlockAndTransactions(height)
		}
	}
}

func saveBlockAndTransactions(height int64) {
	defer func() {
		blockHeightChan <- height
	}()

	entry := logrus.WithFields(logrus.Fields{
		"event":        models.LogEventSaveBlockAndTransactions,
		"block_height": height,
	})

	bestBlock := height < 0
	block, err := wallet.GetBlock(bestBlock, height)
	if err == jerrors.ErrNoNewBlock {
		time.Sleep(time.Minute)
		entry.Info("no new block ahead")
		return
	}

	if err != nil {
		entry.WithError(err).Error("fail to get block from blockchain")
		return
	}

	transactions, err := wallet.GetReceivedSince(block.PrevHash, config.Wallet.MinConfirms)
	if err != nil {
		entry.WithField("hash", block.PrevHash).WithError(err).Error("fail to list transactions from blockchain")
		return
	}

	if err := storage.SaveBlockAndTransactions(walletBlockToModelBlock(block), walletTxsToModelTxs(transactions)); err != nil {
		entry.WithField("hash", block.PrevHash).WithError(err).Error("fail to save block and transactions to db")
		return
	}

	entry.Info("save block and transactions successfully")
	height++
}

func walletTxsToModelTxs(txs []w.Transaction) []models.Transaction {
	transactions := make([]models.Transaction, len(txs))
	for i, v := range txs {
		transactions[i] = models.Transaction{
			Address:        v.Address,
			Amount:         v.Amount,
			TransactionID:  v.TransactionID,
			Hash:           v.Hash,
			BlockCreatedAt: v.BlockCreatedAt,
		}
	}
	return transactions
}

func walletBlockToModelBlock(blockchainBlock *w.Block) models.Block {
	return models.Block{
		Hash:           blockchainBlock.Hash,
		Height:         blockchainBlock.Height,
		BlockCreatedAt: blockchainBlock.BlockCreatedAt.UTC(),
	}
}
