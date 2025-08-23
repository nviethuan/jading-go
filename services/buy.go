package services

import (
	"context"
	"log"
	"strconv"

	"github.com/nviethuan/jading-go/utils"
	binance "github.com/binance/binance-connector-go"
)

func Buy() {
	client := utils.NewBinanceAPI()
	wsApiClient := utils.NewBinanceWSClientAPI()

	err := wsApiClient.Connect()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer wsApiClient.Close()

	depth, err := wsApiClient.NewDepthService().Symbol("BTCUSDT").Limit(10).Do(context.Background())
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	priceStr := depth.Result.Asks[0][0]
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("Error parsing price: %v", err)
		return
	}

	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return
	}

	log.Printf("Account before buy: %v", account)

	order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    Side("BUY").
    Type("LIMIT").
    Quantity(0.0001).
    Price(price).
    TimeInForce("IOC").
    Do(context.Background(), binance.WithRecvWindow(10000))

	if err != nil {
		log.Printf("Error creating order: %v", err)
		return
	}

	log.Printf("Order created: %v", order)

	account, err = client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return
	}

	log.Printf("Account after buy: %v", account)
}
