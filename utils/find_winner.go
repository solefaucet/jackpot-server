package utils

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/jackpot-server/models"
)

// FindWinner finds out the winner address from transactions and block hash
func FindWinner(transactions []models.Transaction, hash string) (string, error) {
	m := map[string]int64{}
	var totalAmount int64
	for _, tx := range transactions {
		amount := int64(tx.Amount * 1e8)
		m[tx.Address] = m[tx.Address] + amount
		totalAmount += amount
	}

	addresses := []string{}
	for address := range m {
		addresses = append(addresses, address)
	}

	sort.Strings(addresses)

	hashLast16Bit := hash[len(hash)-16:]
	randomNumberFromBlockchain, err := strconv.ParseUint(hashLast16Bit, 16, 64)
	if err != nil {
		return "", fmt.Errorf("fail to parse %v to int64: %#v", hashLast16Bit, err)
	}

	sum := int64(randomNumberFromBlockchain%uint64(totalAmount) + 1)

	for _, address := range addresses {
		sum -= m[address]
		if sum <= 0 {
			return address, nil
		}
	}

	// code can never run here
	logrus.WithFields(logrus.Fields{
		"event":         "figure out winner",
		"hash":          hash,
		"hashLast16Bit": hashLast16Bit,
		"random_number": randomNumberFromBlockchain,
		"addresses":     addresses,
		"transactions":  m,
		"sum":           sum,
		"total_amount":  totalAmount,
	}).Panicln("cannot figure out which winner is")
	return "", nil
}
