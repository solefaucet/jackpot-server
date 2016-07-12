package utils

import (
	"testing"

	"github.com/solefaucet/jackpot-server/models"
)

func TestFindWinner(t *testing.T) {
	txs := []models.Transaction{
		{Address: "DCs8E9Gb3mgEweCLFCAuibncGN84znNczs", Amount: 45},
		{Address: "DNNn3syd3RBpRtdA6T1qA7YKU31MoS3whp", Amount: 100},
	}
	expected := "DCs8E9Gb3mgEweCLFCAuibncGN84znNczs"
	actual := FindWinner(txs, "dac75c7c6847bf3e6b983ffed9327970665077c25bc7e69a575e1aefa0aa8d61")

	if actual != expected {
		t.Errorf("address should be %v but get %v", expected, actual)
	}
}
