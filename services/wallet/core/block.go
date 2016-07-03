package core

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/solefaucet/jackpot-server/jerrors"
	"github.com/solefaucet/jackpot-server/services/wallet"
)

// GetBlock get block
func (w Wallet) GetBlock(bestBlock bool, height int64) (*wallet.Block, error) {
	h, err := w.client.GetBlockCount()
	if err != nil {
		return nil, err
	}

	if bestBlock {
		height = h
	}

	if height > h {
		return nil, jerrors.ErrNoNewBlock
	}

	hash, err := w.client.GetBlockHash(height)
	return w.getBlockFromHash(height, hash, err)
}

func (w Wallet) getBlockFromHash(height int64, hash *wire.ShaHash, err error) (*wallet.Block, error) {
	if err != nil {
		return nil, err
	}

	block, err := w.client.GetBlock(hash)
	if err != nil {
		return nil, err
	}

	return &wallet.Block{
		Height:         height,
		Hash:           block.Sha().String(),
		BlockCreatedAt: block.MsgBlock().Header.Timestamp,
	}, nil
}
