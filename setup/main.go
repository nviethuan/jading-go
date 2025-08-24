package main

import (
	"os"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/repositories"
)

func main() {
	repositories.NewAccountRepository().Create(&models.Account{
		Symbol:            "UNIUSDT",
		Network:           "Testnet",
		Description:       "UNIUSDT",
		Email:             "test-uni@test.com",
		ApiKey:            os.Getenv("BINANCE_API_KEY"),
		ApiSecret:         os.Getenv("BINANCE_SECRET_KEY"),
		RestApi:           "https://testnet.binance.vision",
		WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
		WsStream:          "wss://stream.testnet.binance.vision",
		Base:              "UNI",
		Quote:             "USDT",
		Fee:               0.001,
		BaseBalance:       0,
		QuoteBalance:      0,
		Profit:            1,
		IsActived:         1,
		BuyQuantity:       0,
		MaxWithdraw:       100,
		InitialInvestment: 100,
		StepSize:          2,
	})
}
