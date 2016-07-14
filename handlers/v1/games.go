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
	Hash            string             `json:"hash"`
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
	Limit  int64 `form:"limit" binding:"required,min=1,max=10"`
	Offset int64 `form:"offset" binding:"omitempty,min=0"`
}

// Games handler
func Games(
	getGames dependencyGetGames,
	getTransactionsByGameOfs dependencyGetTransactionsByGameOfs,
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
		jackpotAmount := getCurrentJackpotAmount(getGames, fee)

		// get games and transactions
		games, err := getGames(p.Limit, p.Offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		transactions, err := getTransactionsByGameOfs(gameOfs(games)...)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// parse result
		transactionMap := constructTransactionMap(transactions)
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

func gameOfs(games []models.Game) []time.Time {
	gameOfs := make([]time.Time, len(games))
	for i := range games {
		gameOfs[i] = games[i].GameOf
	}
	return gameOfs
}

func getCurrentJackpotAmount(getGames dependencyGetGames, fee float64) float64 {
	games, err := getGames(1, 0)
	if err != nil {
		return 0.0
	}

	if len(games) > 0 {
		return games[0].TotalAmount * (1 - fee)
	}

	return 0.0
}

func constructTransactionMap(transactions []models.Transaction) map[time.Time]map[string]*record {
	transactionMap := make(map[time.Time]map[string]*record)
	for _, v := range transactions {
		t := v.GameOf
		if _, ok := transactionMap[t]; !ok {
			transactionMap[t] = make(map[string]*record)
		}
		if _, ok := transactionMap[t][v.Address]; !ok {
			transactionMap[t][v.Address] = &record{
				Confirmations: v.Confirmations,
				ReceivedAt:    v.BlockCreatedAt,
			}
		}
		transactionMap[t][v.Address].Amount += v.Amount
	}
	return transactionMap
}

func calculateWinProbability(recordMap map[string]*record, totalAmount float64) map[string]*record {
	for _, r := range recordMap {
		r.WinProbability = r.Amount / totalAmount * 100
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

func constructGamesResponse(games []models.Game, transactionMap map[time.Time]map[string]*record, fee float64, blockchainTxURL string) []gameResponse {
	response := make([]gameResponse, len(games))
	for i, v := range games {
		response[i] = gameResponse{
			GameOf:          v.GameOf,
			PaymentProofURL: paymentProofWithTxID(blockchainTxURL, v.TransactionID),
			WinnerAddress:   v.Address,
			Hash:            v.Hash,
			JackpotAmount:   v.TotalAmount * (1 - fee),
			Records:         calculateWinProbability(transactionMap[v.GameOf], v.TotalAmount),
		}
	}
	return response
}
