package services

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	binance "github.com/binance/binance-connector-go"
	"github.com/nviethuan/jading-go/repositories"
	"github.com/nviethuan/jading-go/utils"
)

func start(symbol string, network string, bids *[]binance.Bid, asks *[]binance.Ask) func() {
	return func() {
		account := repositories.NewAccountRepository().FindBySymbol(symbol, network)
		if account == nil {
			fmt.Println("Account not found")
			return
		}

		if account.IsActived == 0 {
			fmt.Println("Account is not actived")
			return
		}

		binanceClient := &utils.Binance{}
		binanceClient.NewBinanceAPI(account.ApiKey, account.ApiSecret, account.RestApi)

		// now := time.Now()
		accountInfo := binanceClient.AccountInfo()
		candlestickData := binanceClient.CandlestickData(symbol, "15m")
		// fmt.Println("Time before: ", time.Since(now))

		accountInfoResponse := <-accountInfo
		candles := <-candlestickData
		// fmt.Println("Time end: ", time.Since(now))

		fmt.Println("accountInfoResponse: ", accountInfoResponse)

		done := make(chan bool, 2)
		var usdtBalance float64 = 0
		var baseBalance float64 = 0

		go func() {
			for _, balance := range accountInfoResponse.Balances {
				if balance.Asset == "USDT" {
					usdtBalance, _ = strconv.ParseFloat(balance.Free, 64)
					done <- true
					break
				}
			}
			done <- true
		}()

		go func() {
			for _, balance := range accountInfoResponse.Balances {
				if balance.Asset == account.Base {
					baseBalance, _ = strconv.ParseFloat(balance.Free, 64)
					done <- true
					break
				}
			}
			done <- true
		}()

		<-done
		<-done

		if usdtBalance < 8 || account.BaseBalance == 0 {
			fmt.Println("USDT balance is less than 8 or base balance is 0")
			return
		}

		oldestCandle := candles[0]
		ask := (*asks)[0]

		oldPrice, _ := strconv.ParseFloat(oldestCandle.Open, 64)
		askPrice, _ := strconv.ParseFloat(ask.Price, 64)
		askValue, _ := strconv.ParseFloat(ask.Quantity, 64)

		isDownTrend := (oldPrice-askPrice)*100/oldPrice >= 1.5
		onSellValue := askValue * askPrice

		// is down trend and usdt balance is less than 70% of on sell value
		if isDownTrend && usdtBalance < onSellValue*0.7 {
			quantity := usdtBalance / askPrice
			quantity = utils.FloorTo(quantity, int(account.StepSize))
			binanceClient.Buy(symbol, quantity, askPrice, "BUY")

			account.BaseBalance = baseBalance + quantity
			account.QuoteBalance = 0
			account.BuyQuantity = askPrice

			repositories.NewAccountRepository().Update(*account)
		}

		candlesLog := [][]any{}
		for _, candle := range candles {
			candlesLog = append(candlesLog, []any{
				time.Unix(int64(candle.OpenTime/1000), 0),
				candle.Low,
				candle.Open,
				candle.Close,
				candle.Open > candle.Close,
			})
		}

		fmt.Println("candles: ", candlesLog)

		fmt.Println("bids: ", *bids)
		fmt.Println("asks: ", *asks)
	}
}

func wsDepthHandler(symbol string, network string) func(event *binance.WsPartialDepthEvent) {
	bids := []binance.Bid{}
	asks := []binance.Ask{}
	throttled := utils.Throttle(start(symbol, network, &bids, &asks), 3*time.Second)
	return func(event *binance.WsPartialDepthEvent) {
		bids = event.Bids
		asks = event.Asks
		throttled()
	}
}

func errHandler(err error) {
	fmt.Println(err)
}

func Run() {

	symbol := flag.String("symbol", "UNIUSDT", "symbol")
	network := flag.String("network", "Testnet", "network")
	flag.Parse()

	websocketStreamClient := binance.NewWebsocketStreamClient(false)

	doneCh, _, err := websocketStreamClient.WsPartialDepthServe100Ms(*symbol, "5", wsDepthHandler(*symbol, *network), errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// // use stopCh to exit
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	stopCh <- struct{}{}
	// }()
	// // remove this if you do not want to be blocked here
	<-doneCh
}
