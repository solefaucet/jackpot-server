package core

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/solefaucet/jackpot-server/services/wallet"
)

// GetBestBlockHash get best block hash
func (w Wallet) GetBestBlockHash() (string, error) {
	return shaHashToWalletHash(w.client.GetBestBlockHash())
}

// GetBlockHash get block hash
func (w Wallet) GetBlockHash(height int64) (string, error) {
	return shaHashToWalletHash(w.client.GetBlockHash(height))
}

func shaHashToWalletHash(hash *wire.ShaHash, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

// GetBlock get block
func (w Wallet) GetBlock(hash string) (wallet.Block, error) {
	shaHash, err := wire.NewShaHashFromStr(hash)
	if err != nil {
		return wallet.Block{}, err
	}

	block, err := w.client.GetBlock(shaHash)
	if err != nil {
		return wallet.Block{}, err
	}

	return wallet.Block{
		Height:         int64(block.Height()),
		Hash:           block.Sha().String(),
		BlockCreatedAt: block.MsgBlock().Header.Timestamp,
	}, nil
}
