package utils

import (
	"context"
	"strconv"

	binance "github.com/binance/binance-connector-go"
)

type Binance struct {
	client                *binance.Client
	websocketAPIClient    *binance.WebsocketAPIClient
	websocketStreamClient *binance.WebsocketStreamClient
}

func (b *Binance) NewBinanceWSClientAPI(apiKey string, secretKey string, baseUrl string) *Binance {
	if b.websocketAPIClient == nil {
		b.websocketAPIClient = binance.NewWebsocketAPIClient(apiKey, secretKey, baseUrl)
	}
	return b
}

func (b *Binance) NewBinanceAPI(apiKey string, secretKey string, baseUrl string) *Binance {
	if b.client == nil {
		b.client = binance.NewClient(apiKey, secretKey, baseUrl)
	}
	return b
}

func (b *Binance) NewBinanceStreamClient() *Binance {
	if b.websocketStreamClient == nil {
		b.websocketStreamClient = binance.NewWebsocketStreamClient(false)
	}
	return b
}

func (b *Binance) AccountInfo() chan binance.AccountResponse {
	accountInfo := make(chan binance.AccountResponse)
	go func() {
		defer close(accountInfo)

		response, err := b.client.NewGetAccountService().Do(context.Background(), binance.WithRecvWindow(10000))
		if err != nil {
			accountInfo <- binance.AccountResponse{}
			return
		}

		var balances []binance.Balance

		// start := time.Now()
		for _, balance := range response.Balances {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			total := free + locked
			if total > 0 {
				balances = append(balances, balance)
			}
		}

		response.Balances = balances

		accountInfo <- *response
	}()

	return accountInfo
}

func (b *Binance) CandlestickData(symbol string, interval string) chan []*binance.KlinesResponse {
	candlestickData := make(chan []*binance.KlinesResponse)
	go func() {
		defer close(candlestickData)

		response, err := b.client.NewKlinesService().Symbol(symbol).Interval(interval).Limit(5).Do(context.Background())
		if err != nil {
			candlestickData <- nil
			return
		}

		candlestickData <- response
	}()

	return candlestickData
}

func (b *Binance) SymbolPriceTicker(symbol string) chan []*binance.TickerPriceResponse {
	symbolPriceTicker := make(chan []*binance.TickerPriceResponse)
	go func() {
		defer close(symbolPriceTicker)

		response, err := b.client.NewTickerPriceService().Symbol(symbol).Do(context.Background())
		if err != nil {
			symbolPriceTicker <- nil
			return
		}

		symbolPriceTicker <- response
	}()

	return symbolPriceTicker
}

func (b *Binance) order(symbol string, side string, quantity float64, price float64, t string) chan binance.CreateOrderResponseRESULT {
	order := make(chan binance.CreateOrderResponseRESULT)
	go func() {
		defer close(order)

		response, err := b.client.NewCreateOrderService().
			Symbol(symbol).
			Side(side).
			Type(t).
			Quantity(quantity).
			Price(price).
			TimeInForce("GTC").
			Do(context.Background(), binance.WithRecvWindow(10000))
		if err != nil {
			order <- binance.CreateOrderResponseRESULT{}
			return
		}

		if orderResponse, ok := response.(binance.CreateOrderResponseRESULT); ok {
			order <- orderResponse
		} else {
			order <- binance.CreateOrderResponseRESULT{}
		}
	}()

	return order
}

func (b *Binance) Buy(symbol string, quantity float64, price float64, t string) chan binance.CreateOrderResponseRESULT {
	return b.order(symbol, "BUY", quantity, price, t)
}

func (b *Binance) Sell(symbol string, quantity float64, price float64, t string) chan binance.CreateOrderResponseRESULT {
	return b.order(symbol, "SELL", quantity, price, t)
}
