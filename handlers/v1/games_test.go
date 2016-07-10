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
		handler := Games(nil, nil, "", time.Minute, 0, "", "", "")

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
		getGamesWithin := mockDependencyGetGamesWithin(nil, fmt.Errorf(""))
		handler := Games(getGamesWithin, nil, "", time.Minute, 0, "", "", "")

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
		getGamesWithin := mockDependencyGetGamesWithin(nil, nil)
		getTransactionsWithin := mockDependencyGetTransactionsWithin(nil, fmt.Errorf(""))
		handler := Games(getGamesWithin, getTransactionsWithin, "", time.Minute, 0, "", "", "")

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
		getGamesWithin := mockDependencyGetGamesWithin(nil, nil)
		getTransactionsWithin := mockDependencyGetTransactionsWithin(nil, nil)
		handler := Games(getGamesWithin, getTransactionsWithin, "", time.Minute, 0, "", "", "")

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
		getGamesWithin := mockDependencyGetGamesWithin(nil, expected)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGamesWithin, time.Now(), time.Now(), 0)

			Convey("Amount should equal 0", func() {
				So(amount, ShouldEqual, 0)
			})
		})
	})

	Convey("Given get games within returning one result", t, func() {
		getGamesWithin := mockDependencyGetGamesWithin([]models.Game{
			{TotalAmount: 100},
		}, nil)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGamesWithin, time.Now(), time.Now(), 0.5)

			Convey("Amount should be 50", func() {
				So(amount, ShouldEqual, 50)
			})
		})
	})

	Convey("Given get games within returning no result", t, func() {
		getGamesWithin := mockDependencyGetGamesWithin([]models.Game{}, nil)

		Convey("When get current jackpot amount", func() {
			amount := getCurrentJackpotAmount(getGamesWithin, time.Now(), time.Now(), 0.5)

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
		{Address: "a1", Amount: 1, BlockCreatedAt: now.Add(3 * time.Minute)},
		{Address: "a2", Amount: 1, BlockCreatedAt: now.Add(2 * time.Minute)},
		{Address: "a1", Amount: 1, BlockCreatedAt: now.Add(time.Minute)},
		{Address: "b1", Amount: 1, BlockCreatedAt: durationAgo.Add(time.Minute)},
		{Address: "b2", Amount: 1, BlockCreatedAt: durationAgo.Add(time.Minute)},
	}

	nowTimestamp := now.Unix()
	durationAgoTimestamp := durationAgo.Unix()
	actual := constructTransactionMap(transactions, duration)
	expected := map[int64]map[string]*record{
		nowTimestamp: map[string]*record{
			"a1": &record{Amount: 2, ReceivedAt: now.Add(3 * time.Minute)},
			"a2": &record{Amount: 1, ReceivedAt: now.Add(2 * time.Minute)},
		},
		durationAgoTimestamp: map[string]*record{
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

	expected := 0.0999
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
	nowTimestamp := now.Unix()
	durationAgoTimestamp := durationAgo.Unix()
	games := []models.Game{
		{TransactionID: "tx_id_1", TotalAmount: 100, GameOf: now},
		{TransactionID: "", TotalAmount: 100, GameOf: durationAgo},
	}
	transactionMap := map[int64]map[string]*record{
		nowTimestamp: map[string]*record{
			"a1": &record{Amount: 2, ReceivedAt: now.Add(3 * time.Minute)},
			"a2": &record{Amount: 1, ReceivedAt: now.Add(2 * time.Minute)},
		},
		durationAgoTimestamp: map[string]*record{
			"b1": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
			"b2": &record{Amount: 1, ReceivedAt: durationAgo.Add(time.Minute)},
		},
	}

	actual := constructGamesResponse(games, transactionMap, 0, "url/")

	transactionMap[nowTimestamp]["a1"].WinProbability = 0.02
	transactionMap[nowTimestamp]["a2"].WinProbability = 0.01
	transactionMap[durationAgoTimestamp]["b1"].WinProbability = 0.01
	transactionMap[durationAgoTimestamp]["b2"].WinProbability = 0.01
	expected := []gameResponse{
		{GameOf: now, JackpotAmount: 100, PaymentProofURL: "url/tx_id_1", Records: transactionMap[nowTimestamp]},
		{GameOf: durationAgo, JackpotAmount: 100, PaymentProofURL: "", Records: transactionMap[durationAgoTimestamp]},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("construct games response expected \n%#v but get \n%#v", expected, actual)
	}
}
