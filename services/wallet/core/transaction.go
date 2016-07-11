package core

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/solefaucet/jackpot-server/models"
	"github.com/solefaucet/jackpot-server/services/wallet"
)

// GetReceivedSince returns transactions since block
func (w Wallet) GetReceivedSince(prevHash, curHash string) ([]wallet.Transaction, error) {
	blockHash, err := wire.NewShaHashFromStr(prevHash)
	if err != nil {
		return nil, err
	}

	result, err := w.client.ListSinceBlock(blockHash)
	if err != nil {
		return nil, err
	}

	var transactions []wallet.Transaction
	for i := len(result.Transactions) - 1; i >= 0; i-- {
		tx := result.Transactions[i]
		if tx.Category != "receive" || tx.BlockHash != curHash {
			continue
		}

		senderAddress, err := w.getSenderAddress(tx.TxID)
		if err != nil {
			return nil, err
		}

		transaction := wallet.Transaction{
			Address:        senderAddress,
			Amount:         tx.Amount,
			TransactionID:  tx.TxID,
			Hash:           tx.BlockHash,
			Confirmations:  tx.Confirmations,
			BlockCreatedAt: time.Unix(tx.BlockTime, 0),
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// NOTE: getSenderAddress is a bit tricky, not sure if it's reliable or not
func (w Wallet) getSenderAddress(txid string) (string, error) {
	resultVin, err := w.getRawTransactionResult(txid)
	if err != nil {
		return "", err
	}
	vin := resultVin.Vin[0]

	resultVout, err := w.getRawTransactionResult(vin.Txid)
	if err != nil {
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"event":          models.LogEventGetSenderAddress,
		"tx_id":          txid,
		"result_vin":     resultVin,
		"result_vout":    resultVout,
		"vin_tx_id":      vin.Txid,
		"vin_vout_index": vin.Vout,
	}).Debug("get sender address information")

	return resultVout.Vout[int(vin.Vout)].ScriptPubKey.Addresses[0], nil
}

func (w Wallet) getRawTransactionResult(txid string) (*btcjson.TxRawResult, error) {
	entry := logrus.WithFields(logrus.Fields{
		"event": models.LogEventGetRawTransaction,
		"tx_id": txid,
	})

	txHash, err := wire.NewShaHashFromStr(txid)
	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to create sha hash from tx id")
		return nil, err
	}

	result, err := w.client.GetRawTransactionVerbose(txHash)
	if err != nil {
		entry.WithField("error", err.Error()).Error("fail to get raw transaction result")
		return nil, err
	}

	return result, nil
}

// SendFromAccountToAddress send coin to address, return transaction id
func (w Wallet) SendFromAccountToAddress(account, address string, amount float64) (string, error) {
	hash, err := w.client.SendFrom(account, address, amount)
	if err != nil {
		return "", fmt.Errorf("core wallet send to address error: %#v", err)
	}

	return hash.String(), nil
}

// GetConfirmationsFromTxID returns confirmations given tx id
func (w Wallet) GetConfirmationsFromTxID(txid string) (int64, error) {
	result, err := w.getRawTransactionResult(txid)
	if err != nil {
		return 0, err
	}

	return int64(result.Confirmations), nil
}
