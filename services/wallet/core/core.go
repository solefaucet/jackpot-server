package core

import (
	"github.com/btcsuite/btcrpcclient"
	"github.com/solefaucet/jackpot-server/services/wallet"
)

// Wallet implements Wallet interface for blockchain manipulation
type Wallet struct {
	client *btcrpcclient.Client
}

var _ wallet.Wallet = Wallet{}

// New create
func New(rpchost, rpcusername, rpcpassword string) (Wallet, error) {
	config := &btcrpcclient.ConnConfig{
		Host:         rpchost,
		User:         rpcusername,
		Pass:         rpcpassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := btcrpcclient.New(config, nil)

	return Wallet{client: client}, err
}
