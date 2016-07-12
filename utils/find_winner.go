package utils

import (
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/models"
)

// FindWinner finds out the winner address from transactions and block hash
func FindWinner(transactions []models.Transaction, hash string) string {
	totalAmount := totalAmountOfTransactions(transactions)
	m := transactionMap(transactions)
	addresses := sortedAddresses(m)
	sum := randomSum(hash, totalAmount)

	entry := logrus.WithFields(logrus.Fields{
		"event":        "figure out winner",
		"hash":         hash,
		"addresses":    addresses,
		"transactions": m,
		"sum":          sum,
		"total_amount": totalAmount,
	})

	for _, address := range addresses {
		sum -= m[address]
		if sum <= 0 {
			entry.WithField("winner_address", address).Info("winner address found")
			return address
		}
	}

	// code can never run here
	entry.Panicln("cannot figure out which winner is")
	return ""
}

func totalAmountOfTransactions(transactions []models.Transaction) int64 {
	var totalAmount int64
	for _, tx := range transactions {
		totalAmount += int64(tx.Amount * 1e8)
	}
	return totalAmount
}

func transactionMap(transactions []models.Transaction) map[string]int64 {
	m := map[string]int64{}
	for _, tx := range transactions {
		m[tx.Address] = m[tx.Address] + int64(tx.Amount*1e8)
	}
	return m
}

func sortedAddresses(transactionMap map[string]int64) []string {
	addresses := []string{}
	for address := range transactionMap {
		addresses = append(addresses, address)
	}
	sort.Strings(addresses)
	return addresses
}

func randomSum(hash string, totalAmount int64) int64 {
	hashLast16Hexadecimal := hash[len(hash)-16:]
	randomNumberFromBlockchain, _ := strconv.ParseUint(hashLast16Hexadecimal, 16, 64)
	return int64(randomNumberFromBlockchain%uint64(totalAmount) + 1)
}
