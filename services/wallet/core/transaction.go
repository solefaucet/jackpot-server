package core

import (
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/solefaucet/jackpot-server/services/wallet"
)

// GetReceivedSince returns transactions since block and minConfirms
func (w Wallet) GetReceivedSince(hash string, minConfirms int) ([]wallet.Transaction, error) {
	blockHash, err := wire.NewShaHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	result, err := w.client.ListSinceBlockMinConf(blockHash, minConfirms)
	if err != nil {
		return nil, err
	}

	var transactions []wallet.Transaction
	for i := len(result.Transactions) - 1; i >= 0; i-- {
		tx := result.Transactions[i]
		if tx.Category != "receive" {
			continue
		}
		transaction := wallet.Transaction{
			Address:        tx.Address,
			Amount:         tx.Amount,
			TransactionID:  tx.TxID,
			Hash:           tx.BlockHash,
			BlockCreatedAt: time.Unix(tx.BlockTime, 0),
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
