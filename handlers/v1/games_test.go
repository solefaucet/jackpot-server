package v1

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/jackpot-server/models"
)

func TestGames(t *testing.T) {
	Convey("Given games handler", t, func() {
		handler := Games(nil, nil, "", time.Minute, 0, "", "", "", "")

		Convey("When request games handler with incorrect parameter", func() {
			route := "/games"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", "/games?offset=-1", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})

	Convey("Given games handler with errored get games within", t, func() {
		getGames := mockDependencyGetGames(nil, fmt.Errorf(""))
		handler := Games(getGames, nil, "", time.Minute, 0, "", "", "", "")

		Convey("When request games handler", func() {
			route := "/games"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", "/games?offset=1&limit=1", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given games handler with errored get transactions within", t, func() {
		getGames := mockDependencyGetGames(nil, nil)
		getTransactionsWithin := mockDependencyGetTransactionsByGameOfs(nil, fmt.Errorf(""))
		handler := Games(getGames, getTransactionsWithin, "", time.Minute, 0, "", "", "", "")

		Convey("When request games handler", func() {
			route := "/games"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", "/games?offset=1&limit=1", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given games handler with everything correct", t, func() {
		getGames := mockDependencyGetGames([]models.Game{
			{},
		}, nil)
		getTransactionsWithin := mockDependencyGetTransactionsByGameOfs(nil, nil)
		handler := Games(getGames, getTransactionsWithin, "", time.Minute, 0, "", "", "", "")

		Convey("When request games handler", func() {
			route := "/games"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", "/games?offset=1&limit=1", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestGetCurrentJackpotAmount(t *testing.T) {
	Convey("Given errored get games within", t, func() {
		expected := fmt.Errorf("error")
		getGames := mockDependencyGetGames(nil, expected)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGames, 0)

			Convey("Amount should equal 0", func() {
				So(amount, ShouldEqual, 0)
			})
		})
	})

	Convey("Given get games within returning one result", t, func() {
		getGames := mockDependencyGetGames([]models.Game{
			{TotalAmount: 100},
		}, nil)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGames, 0.5)

			Convey("Amount should be 50", func() {
				So(amount, ShouldEqual, 50)
			})
		})
	})

	Convey("Given get games within returning no result", t, func() {
		getGames := mockDependencyGetGames([]models.Game{}, nil)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGames, 0.5)

			Convey("Amount should be 0", func() {
				So(amount, ShouldEqual, 0)
			})
		})
	})
}

func TestConstructTransactionMap(t *testing.T) {
	duration := time.Hour
	now := time.Now().Truncate(duration)
	durationAgo := now.Add(-duration)
	transactions := []models.Transaction{
		{Address: "a1", Amount: 1, GameOf: now, BlockCreatedAt: now.Add(3 * time.Minute)},
		{Address: "a2", Amount: 1, GameOf: now, BlockCreatedAt: now.Add(2 * time.Minute)},
		{Address: "a1", Amount: 1, GameOf: now, BlockCreatedAt: now.Add(time.Minute)},
		{Address: "b1", Amount: 1, GameOf: durationAgo, BlockCreatedAt: durationAgo.Add(time.Minute)},
		{Address: "b2", Amount: 1, GameOf: durationAgo, BlockCreatedAt: durationAgo.Add(time.Minute)},
	}

	actual := constructTransactionMap(transactions)
	expected := map[time.Time]map[string]*record{
		now: map[string]*record{
			"a1": &record{Amount: 2, ReceivedAt: now.Add(3 * time.Minute)},
			"a2": &record{Amount: 1, ReceivedAt: now.Add(2 * time.Minute)},
		},
		durationAgo: map[string]*record{
			"b1": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
			"b2": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("construct transaction map expected \n%#v but get \n%#v", expected, actual)
	}
}

func TestCalculateWinProbability(t *testing.T) {
	recordMap := map[string]*record{
		"key": &record{Amount: 9.99},
	}
	calculateWinProbability(recordMap, 100)

	expected := 9.99
	if actual := recordMap["key"].WinProbability; actual != expected {
		t.Errorf("calculate win probability expected %v but get %v", expected, actual)
	}
}

func TestPaymentProofWithTxID(t *testing.T) {
	if actual, expected := paymentProofWithTxID("url", ""), ""; actual != expected {
		t.Errorf("payment proof with txid expected %v but get %v", expected, actual)
	}

	if actual, expected := paymentProofWithTxID("url", "123"), "url123"; actual != expected {
		t.Errorf("payment proof with txid expected %v but get %v", expected, actual)
	}
}

func TestConstructGamesResponse(t *testing.T) {
	duration := time.Hour
	now := time.Now().Truncate(duration)
	durationAgo := now.Add(-duration)
	games := []models.Game{
		{TransactionID: "tx_id_1", TotalAmount: 100, GameOf: now},
		{TransactionID: "", TotalAmount: 100, GameOf: durationAgo},
	}
	transactionMap := map[time.Time]map[string]*record{
		now: map[string]*record{
			"a1": &record{Amount: 2, ReceivedAt: now.Add(3 * time.Minute)},
			"a2": &record{Amount: 1, ReceivedAt: now.Add(2 * time.Minute)},
		},
		durationAgo: map[string]*record{
			"b1": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
			"b2": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
		},
	}

	actual := constructGamesResponse(games, transactionMap, 0, "url/")

	transactionMap[now]["a1"].WinProbability = 0.02
	transactionMap[now]["a2"].WinProbability = 0.01
	transactionMap[durationAgo]["b1"].WinProbability = 0.01
	transactionMap[durationAgo]["b2"].WinProbability = 0.01
	expected := []gameResponse{
		{GameOf: now, JackpotAmount: 100, PaymentProofURL: "url/tx_id_1", Records: transactionMap[now]},
		{GameOf: durationAgo, JackpotAmount: 100, PaymentProofURL: "", Records: transactionMap[durationAgo]},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("construct games response expected \n%#v but get \n%#v", expected, actual)
	}
}
