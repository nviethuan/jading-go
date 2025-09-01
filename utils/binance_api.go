package utils

import (
	"context"
	"fmt"
	"strconv"

	binance "github.com/binance/binance-connector-go"
	"github.com/nviethuan/jading-go/models"
)

type Binance struct {
	client                *binance.Client
	websocketAPIClient    *binance.WebsocketAPIClient
	websocketStreamClient *binance.WebsocketStreamClient
}

func (b *Binance) NewBinanceWSClientAPI(account *models.Account) *Binance {
	if b.websocketAPIClient == nil {
		b.websocketAPIClient = binance.NewWebsocketAPIClient(account.ApiKey, account.ApiSecret, account.RestApi)
	}
	return b
}

func (b *Binance) NewBinanceAPI(account *models.Account) *Binance {
	if b.client == nil {
		client := binance.NewClient(account.ApiKey, account.ApiSecret, account.RestApi)
		b.client = client
	}
	return b
}

func (b *Binance) NewBinanceStreamClient() *Binance {
	if b.websocketStreamClient == nil {
		b.websocketStreamClient = binance.NewWebsocketStreamClient(false)
	}
	return b
}

func (b *Binance) AccountInfo() chan *binance.AccountResponse {
	accountInfo := make(chan *binance.AccountResponse, 1)
	go func() {
		defer close(accountInfo)

		response, err := b.client.NewGetAccountService().Do(context.Background(), binance.WithRecvWindow(10000))
		if err != nil {
			accountInfo <- nil
			return
		}

		var balances []binance.Balance

		for _, balance := range response.Balances {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			total := free + locked
			if total > 0 {
				balances = append(balances, balance)
			}
		}

		response.Balances = balances

		accountInfo <- response
	}()

	return accountInfo
}

func (b *Binance) TradeFee(ctx context.Context, symbol string) chan []*binance.TradeFeeResponse {
	tradeFee := make(chan []*binance.TradeFeeResponse, 1)
	go func() {
		defer close(tradeFee)

		response, err := b.client.NewTradeFeeService().Symbol(symbol).Do(ctx)
		if err != nil {
			tradeFee <- []*binance.TradeFeeResponse{}
			return
		}

		tradeFee <- response
	}()

	return tradeFee
}

func (b *Binance) AccountInfoWithContext(ctx context.Context) chan binance.AccountResponse {
	accountInfo := make(chan binance.AccountResponse, 1)
	go func() {
		defer close(accountInfo)

		// Tạo một channel để nhận kết quả từ API
		resultCh := make(chan *binance.AccountResponse, 1)
		errCh := make(chan error, 1)

		go func() {
			response, err := b.client.NewGetAccountService().Do(ctx, binance.WithRecvWindow(10000))
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- response
		}()

		select {
		// Khi context bị huỷ (ví dụ: timeout, cancel), thoát goroutine mà không gửi gì vào accountInfo
		case <-ctx.Done():
			// Nếu context bị huỷ, không gửi gì vào accountInfo, chỉ return
			return
		case err := <-errCh:
			accountInfo <- binance.AccountResponse{}
			fmt.Println("[Binance] Error: ", err)
			return
		case response := <-resultCh:
			var balances []binance.Balance
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
		}
	}()

	return accountInfo
}

func (b *Binance) CandlestickData(symbol string, interval string) chan []*binance.KlinesResponse {
	candlestickData := make(chan []*binance.KlinesResponse, 1)
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
	symbolPriceTicker := make(chan []*binance.TickerPriceResponse, 1)
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

func (b *Binance) order(symbol string, side string, quantity float64, price float64, t string) chan binance.CreateOrderResponseFULL {
	order := make(chan binance.CreateOrderResponseFULL, 1)
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
			fmt.Println("[Binance] Error: ", err)
			order <- binance.CreateOrderResponseFULL{}
			return
		}

		orderResponse, ok := response.(*binance.CreateOrderResponseFULL)

		if ok {
			fmt.Println("[Binance] Order Response: ", orderResponse)
			order <- *orderResponse
		} else {
			order <- binance.CreateOrderResponseFULL{}
		}
	}()

	return order
}

func (b *Binance) Buy(symbol string, quantity float64, price float64, t string) chan binance.CreateOrderResponseFULL {
	return b.order(symbol, "BUY", quantity, price, t)
}

func (b *Binance) Sell(symbol string, quantity float64, price float64, t string) chan binance.CreateOrderResponseFULL {
	return b.order(symbol, "SELL", quantity, price, t)
}

func (b *Binance) Withdraw(asset string, quantity float64) chan binance.TransferToMasterResp {
	withdraw := make(chan binance.TransferToMasterResp, 1)
	go func() {
		defer close(withdraw)

		response, err := b.client.NewTransferToMasterService().
			Asset(asset).
			Amount(quantity).
			Do(context.Background(), binance.WithRecvWindow(10000))

		if err != nil {
			withdraw <- binance.TransferToMasterResp{}
			return
		}

		withdraw <- *response
	}()

	return withdraw
}
