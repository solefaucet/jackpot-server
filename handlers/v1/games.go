package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/jackpot-server/models"
)

type gamesResponse struct {
	DestAddress    string         `json:"dest_address"`
	DestAddressURL string         `json:"dest_address_url"`
	Duration       int64          `json:"duration"`
	QRCode         string         `json:"qrcode"`
	JackpotAmount  float64        `json:"jackpot_amout"`
	NextGameTime   time.Time      `json:"next_game_time"`
	Games          []gameResponse `json:"games"`
}

type gameResponse struct {
	GameOf          time.Time          `json:"game_of"`
	PaymentProofURL string             `json:"payment_proof_url"`
	WinnerAddress   string             `json:"winner_address"`
	JackpotAmount   float64            `json:"jackpot_amount"`
	Records         map[string]*record `json:"records"`
}

type record struct {
	Confirmations  int64     `json:"confirmations"`
	Amount         float64   `json:"amount"`
	WinProbability float64   `json:"win_probability"`
	ReceivedAt     time.Time `json:"received_at"`
}

type gamePayload struct {
	Limit  time.Duration `form:"limit" binding:"required,min=1,max=10"`
	Offset time.Duration `form:"offset" binding:"omitempty,min=0"`
}

// Games handler
func Games(
	getGamesWithin dependencyGetGamesWithin,
	getTransactionsWithin dependencyGetTransactionsWithin,
	destAddress string,
	duration time.Duration,
	fee float64,
	blockchainTxURL, blockchainAddressURL, coinType, label string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := gamePayload{}
		if err := c.BindWith(&p, binding.Form); err != nil {
			return
		}

		// get current jackpot amount
		now := time.Now()
		jackpotAmount := getCurrentJackpotAmount(getGamesWithin, now.Truncate(duration), now.Truncate(duration).Add(duration), fee)

		// get games and transactions
		end := now.Truncate(duration).Add(duration).Add(-duration * p.Offset)
		start := end.Add(-duration * p.Limit)

		games, err := getGamesWithin(start, end)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		transactions, err := getTransactionsWithin(start, end)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// parse result
		transactionMap := constructTransactionMap(transactions, duration)
		gs := constructGamesResponse(games, transactionMap, fee, blockchainTxURL)

		response := gamesResponse{
			Games:          gs,
			Duration:       duration.Nanoseconds() / 1e9,
			DestAddress:    destAddress,
			DestAddressURL: blockchainAddressURL + destAddress,
			JackpotAmount:  jackpotAmount,
			NextGameTime:   now.Truncate(duration).Add(duration),
			QRCode:         fmt.Sprintf("%s:%s?label=%s", coinType, destAddress, label),
		}

		// response
		c.JSON(http.StatusOK, response)
	}
}

func getCurrentJackpotAmount(getGamesWithin dependencyGetGamesWithin, start, end time.Time, fee float64) float64 {
	games, err := getGamesWithin(start, end)
	if err != nil {
		return 0.0
	}

	if len(games) > 0 {
		return games[0].TotalAmount * (1 - fee)
	}

	return 0.0
}

func constructTransactionMap(transactions []models.Transaction, duration time.Duration) map[int64]map[string]*record {
	transactionMap := make(map[int64]map[string]*record)
	for _, v := range transactions {
		timestamp := v.BlockCreatedAt.Truncate(duration).Unix()
		if _, ok := transactionMap[timestamp]; !ok {
			transactionMap[timestamp] = make(map[string]*record)
		}
		if _, ok := transactionMap[timestamp][v.Address]; !ok {
			transactionMap[timestamp][v.Address] = &record{
				Confirmations: v.Confirmations,
				ReceivedAt:    v.BlockCreatedAt,
			}
		}
		transactionMap[timestamp][v.Address].Amount += v.Amount
	}
	return transactionMap
}

func calculateWinProbability(recordMap map[string]*record, totalAmount float64) map[string]*record {
	for _, r := range recordMap {
		r.WinProbability = r.Amount / totalAmount
	}
	return recordMap
}

func paymentProofWithTxID(url, txid string) string {
	paymentProofURL := ""
	if txid != "" {
		paymentProofURL = url + txid
	}
	return paymentProofURL
}

func constructGamesResponse(games []models.Game, transactionMap map[int64]map[string]*record, fee float64, blockchainTxURL string) []gameResponse {
	response := make([]gameResponse, len(games))
	for i, v := range games {
		response[i] = gameResponse{
			GameOf:          v.GameOf,
			PaymentProofURL: paymentProofWithTxID(blockchainTxURL, v.TransactionID),
			WinnerAddress:   v.Address,
			JackpotAmount:   v.TotalAmount * (1 - fee),
			Records:         calculateWinProbability(transactionMap[v.GameOf.Unix()], v.TotalAmount),
		}
	}
	return response
}
