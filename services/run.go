package services

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	binance "github.com/binance/binance-connector-go"
	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/repositories"
	"github.com/nviethuan/jading-go/utils"
	"github.com/slack-go/slack"
)

var slackClient *utils.SlackClient
var binanceClient *utils.Binance

func processBuy(account *models.Account, symbol *string, asks *[]binance.Ask, usdtBalance float64, candles *[]*binance.KlinesResponse) {
	// check if usdt balance is less than 8 and base balance is 0
	// stop if we have sold all base balance and usdt balance is less than 8
	hasSold := account.BaseBalance == 0
	if usdtBalance < account.MinStopLoss && !hasSold {
		fmt.Println("[", symbol, "] STOP! USDT balance is less than 8 and base balance is 0")
		return
	}
	// ------------

	// check downtrend
	oldestCandle := (*candles)[0]
	// get first sell order
	ask := (*asks)[0]

	oldPrice, _ := strconv.ParseFloat(oldestCandle.Open, 64)
	askPrice, _ := strconv.ParseFloat(ask.Price, 64)
	askValue, _ := strconv.ParseFloat(ask.Quantity, 64)

	// take the open price of the oldest candle and compare it with the best ask price
	// if the ask price drops more than 1.5% below the open price of the oldest candle, consider it a downtrend
	isDownTrend := (oldPrice-askPrice)/oldPrice >= account.Threshold // ✅
	// ------------

	// check if the 70% of the quantity they want to sell
	// is available based on the amount we bought
	onSellValue := (askValue * askPrice) * 0.7
	quantity := usdtBalance / askPrice
	isEnoughUsdtBalance := quantity <= onSellValue
	// ------------

	// combine all conditions
	if isDownTrend && isEnoughUsdtBalance {
		// calculate the quantity we want to buy
		quantity = utils.FloorTo(quantity, int(account.StepSize))
		// ------

		// buy #########################################################
		buyChan := binanceClient.Buy(account.Symbol, quantity, askPrice, "BUY")
		buyResponse := <-buyChan
		// --- #########################################################

		// update state to database
		account.BaseBalance = account.BaseBalance + quantity
		account.QuoteBalance = 0
		account.BuyQuantity = usdtBalance

		repositories.NewAccountRepository().Update(*account)
		// ------------

		// log to slack
		title := fmt.Sprintf("💰 Buy %f (%s) with %f", quantity, strings.ToUpper(account.Base), askPrice)
		msg := fmt.Sprintf(":%s: :dollar: [BUY] %f (%s) with price *%f* - order id: `%d`",
			strings.ToLower(account.Symbol), // emoji
			quantity,
			strings.ToUpper(account.Base),
			askPrice,
			buyResponse.OrderId,
		)

		bodyText := slack.NewTextBlockObject("mrkdwn", msg, false, true)
		bodyBlock := slack.NewSectionBlock(bodyText, nil, nil)

		blocks := []slack.Block{bodyBlock}

		tsChan := slackClient.SendInfo(title, "", blocks...)
		// ------------

		// create stack trade
		ts := <-tsChan

		priceSell := usdtBalance / (quantity * (1 - account.Fee) * (1 + account.Profit))

		now := time.Now()

		// (price_sell <= ? OR price_buy >= ?)
		repositories.NewStackTradeRepository().Create(models.StackTrade{
			Symbol:    account.Symbol,
			Quantity:  quantity,
			PriceBuy:  askPrice,
			PriceSell: priceSell,
			ThreadID:  ts,
			Status:    "BUY",
			CreatedAt: now,
			UpdatedAt: now,
		})
		// ------------

		return
	}
}

func processSell(account *models.Account, bids *[]binance.Bid, usdtBalance float64) {
	for i := len(*bids) - 1; i >= 0; i-- {
		bid := (*bids)[i]
		bidPrice, _ := strconv.ParseFloat(bid.Price, 64)
		bidValue, _ := strconv.ParseFloat(bid.Quantity, 64)

		stopLoss := bidPrice * (1 - account.StopLoss)

		stackTrades := repositories.NewStackTradeRepository().FindBySymbol(account.Symbol, "BUY", bidPrice, bidValue*0.07, stopLoss)

		if len(stackTrades) == 0 {
			continue
		}

		stackTrade := stackTrades[0]
		isStopLoss := stackTrade.PriceBuy >= stopLoss

		purpose := "*sell*"
		if isStopLoss {
			purpose = "`stoploss`"
		}

		sellChan := binanceClient.Sell(account.Symbol, stackTrade.Quantity, bidPrice, "LIMIT")
		sellResponse := <-sellChan

		quantityEarn, _ := strconv.ParseFloat(sellResponse.OrigQty, 64)

		shouldWithdraw := usdtBalance+quantityEarn > account.MaxWithdraw
		withdrawQuantity := account.MaxWithdraw - (usdtBalance + quantityEarn)

		if shouldWithdraw {
			<-binanceClient.Withdraw("USDT", withdrawQuantity)
		}

		stackTrade.Status = "SELL"
		stackTrade.PriceSell = bidPrice
		stackTrade.UpdatedAt = time.Now()
		repositories.NewStackTradeRepository().Update(*stackTrade)

		// quote balance = quantity has bought * current price sell
		quoteBalance := stackTrade.Quantity * bidPrice

		account.BaseBalance = math.Max(0, account.BaseBalance-stackTrade.Quantity)
		account.QuoteBalance = quoteBalance

		repositories.NewAccountRepository().Update(*account)

		// log to slack
		title := fmt.Sprintf("💰 Sell %f (%s) with %f", stackTrade.Quantity, strings.ToUpper(account.Base), bidPrice)
		msg := fmt.Sprintf(":%s: :dollar: [SELL] %f (%s) with price *%f* - order id: `%d`\nby %s",
			strings.ToLower(account.Symbol), // emoji
			stackTrade.Quantity,
			strings.ToUpper(account.Base),
			bidPrice,
			stackTrade.ID,
			purpose,
		)

		shouldWithdrawMsg := fmt.Sprintf("💰:%s: *No withdraw*: `%f` (USDT)", strings.ToUpper(account.Symbol), quoteBalance)

		if shouldWithdraw {
			shouldWithdrawMsg = fmt.Sprintf("💰:%s: *Withdraw*: `%f` (USDT)", strings.ToUpper(account.Symbol), withdrawQuantity)
		}

		bodyText := slack.NewTextBlockObject("mrkdwn", msg, false, true)
		shouldWithdrawText := slack.NewTextBlockObject("mrkdwn", shouldWithdrawMsg, false, true)

		bodyBlock1 := slack.NewSectionBlock(bodyText, nil, nil)
		bodyBlock2 := slack.NewSectionBlock(shouldWithdrawText, nil, nil)

		blocks := []slack.Block{bodyBlock1}

		if shouldWithdraw {
			blocks = append(blocks, bodyBlock2)
		}

		<-slackClient.SendInfo(title, "", blocks...)

		return
	}
}

func start(symbol *string, network *string, bids *[]binance.Bid, asks *[]binance.Ask) func() {
	return func() {
		// get account from database
		account := repositories.NewAccountRepository().FindBySymbol(symbol, network)
		if account == nil {
			fmt.Println("[", *symbol, "] Account not found")
			return
		}

		if account.IsActived == 0 {
			fmt.Println("[", symbol, "] Account is not actived")
			return
		}
		// -----------

		// init binance client with account
		binanceClient.NewBinanceAPI(account)
		// ------------

		// get account info and candlestick data
		accountInfoChan := binanceClient.AccountInfo()
		candlestickDataChan := binanceClient.CandlestickData(*symbol, "15m")

		accountInfo := <-accountInfoChan
		candles := <-candlestickDataChan
		// ------------

		// take my usdt balance and my base balance
		usdtBalance := 0.0

		// get usdt balance
		for _, balance := range accountInfo.Balances {
			if balance.Asset == "USDT" {
				num, _ := strconv.ParseFloat(balance.Free, 64)
				usdtBalance = num
				break
			}
		}
		// ------------

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			processBuy(account, symbol, asks, usdtBalance, &candles)
		}()
		go func() {
			defer wg.Done()
			processSell(account, bids, usdtBalance)
		}()
		wg.Wait()
	}
}

func wsDepthHandler(symbol *string, network *string) func(event *binance.WsPartialDepthEvent) {
	bids := []binance.Bid{}
	asks := []binance.Ask{}
	throttled := utils.Throttle(start(symbol, network, &bids, &asks), 3*time.Second)
	return func(event *binance.WsPartialDepthEvent) {
		bids = event.Bids
		asks = event.Asks
		throttled()
	}
}

func startTest(bids *[]binance.Bid, asks *[]binance.Ask) func() {
	return func() {
		fmt.Println("bids", binance.PrettyPrint(bids))
		fmt.Println("asks", binance.PrettyPrint(asks))
	}
}

func wsDepthHandlerTest() func(event *binance.WsPartialDepthEvent) {
	bids := []binance.Bid{}
	asks := []binance.Ask{}

	throttled := utils.Throttle(startTest(&bids, &asks), 3*time.Second)

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
	slackClient = &utils.SlackClient{}
	slackClient.NewSlackClient()
	symbolStr := flag.String("symbol", "UNIUSDT", "symbol")
	network := flag.String("network", "Testnet", "network")
	flag.Parse()

	websocketStreamClient := binance.NewWebsocketStreamClient(false, "wss://stream.testnet.binance.vision")

	symbols := strings.Split(*symbolStr, ",")

	doneChs := make([]chan struct{}, len(symbols))
	wg := sync.WaitGroup{}

	for i, symbol := range symbols {
		wg.Add(1)
		go func(idx int, sym string) {
			defer wg.Done()
			doneCh, _, err := websocketStreamClient.WsPartialDepthServe100Ms(sym, "10", wsDepthHandler(&sym, network), errHandler)
			if err != nil {
				fmt.Println(err)
				return
			}
			doneChs[idx] = doneCh
		}(i, symbol)
	}

	wg.Wait()

	for _, doneCh := range doneChs {
		if doneCh != nil {
			<-doneCh
		}
	}
}

func RunTest() {
	slackClient = &utils.SlackClient{}
	slackClient.NewSlackClient()
	symbolStr := flag.String("symbol", "UNIUSDT", "symbol")
	flag.Parse()

	websocketStreamClient := binance.NewWebsocketStreamClient(false)

	symbols := strings.Split(*symbolStr, ",")

	doneChs := make([]chan struct{}, len(symbols))
	wg := sync.WaitGroup{}

	for i, symbol := range symbols {
		wg.Add(1)
		go func(idx int, sym string) {
			defer wg.Done()
			doneCh, _, err := websocketStreamClient.WsPartialDepthServe100Ms(sym, "10", wsDepthHandlerTest(), errHandler)
			if err != nil {
				fmt.Println(err)
				return
			}
			doneChs[idx] = doneCh
		}(i, symbol)
	}

	wg.Wait()

	for _, doneCh := range doneChs {
		if doneCh != nil {
			<-doneCh
		}
	}
}
