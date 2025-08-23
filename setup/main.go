package main

import (
	"os"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/repositories"
)

func main() {
	repositories.NewAccountRepository().Create(&models.Account{
		Symbol:            "BTCUSDT",
		Network:           "Testnet",
		Description:       "BTCUSDT",
		Email:             "test@test.com",
		ApiKey:            os.Getenv("BINANCE_API_KEY"),
		ApiSecret:         os.Getenv("BINANCE_SECRET_KEY"),
		RestApi:           "https://testnet.binance.vision",
		WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
		WsStream:          "wss://stream.testnet.binance.vision",
		Base:              "BTC",
		Quote:             "USDT",
		Fee:               0.001,
		Profit:            0.001,
		IsActived:         1,
		BuyPrice:          0,
		MaxWithdraw:       100,
		InitialInvestment: 100,
		StepSize:          0.00001,
	})
}
