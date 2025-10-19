package services

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
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
var binanceClient *utils.Binance = &utils.Binance{}

func calculateRSI(candles *[]*binance.KlinesResponse) float64 {
	ups := []float64{}
	downs := []float64{}

	rsi := 0.0

	for i := 1; i < len(*candles); i++ {
		prevClose, _ := strconv.ParseFloat((*candles)[i-1].Close, 64)
		currClose, _ := strconv.ParseFloat((*candles)[i].Close, 64)
		diff := currClose - prevClose

		if diff > 0 {
			ups = append(ups, diff)
		} else {
			downs = append(downs, math.Abs(diff))
		}
	}

	avgUp := utils.Average(ups)
	avgDown := utils.Average(downs)

	rs := avgUp / avgDown
	rsi = 100 - (100 / (1 + rs))

	return rsi
}

func processBuy(t string, account *models.Account, asks *[]binance.Ask, usdtBalance float64, baseBalance float64, loc *time.Location) {
	prefixLog := fmt.Sprintf("%s BUY_%s: ", t, account.Symbol)
	fmt.Printf("%s Process Buy =======\n", prefixLog)
	// check if usdt balance is less than 8 and base balance is 0
	// stop if we have sold all base balance and usdt balance is less than 8
	hasSold := account.BaseBalance == 0
	if usdtBalance < account.MinStopLoss && !hasSold {
		fmt.Printf("%s STOP! USDT balance is less than 8 and base balance is not 0\n", prefixLog)
		return
	}
	// ------------

	candlestickDataChan := binanceClient.CandlestickData(account, account.Symbol, "1m", 15)
	if candlestickDataChan == nil {
		fmt.Printf("%s %s - STOP! CandlestickData is not available\n", t, account.Symbol)
		return
	}

	candles := <-candlestickDataChan

	// check downtrend
	oldestCandle := candles[0]
	// get first sell order
	ask := (*asks)[0]

	oldestPrice, _ := strconv.ParseFloat(oldestCandle.Open, 64)
	oldestTime := time.UnixMilli(int64(oldestCandle.OpenTime)).In(loc).Format("2006-01-02 15:04:05")
	askPrice, _ := strconv.ParseFloat(ask.Price, 64)
	askValue, _ := strconv.ParseFloat(ask.Quantity, 64)

	// take the open price of the oldest candle and compare it with the best ask price
	// if the ask price drops more than the threshold (default: 0.015 ~ 1.5%) below the open price of the oldest candle, consider it a downtrend
	// isDownTrend := (oldestPrice-askPrice)/oldestPrice >= account.Threshold // ✅
	// ------------

	// should remove the newest candle because it's not closed yet
	cds := candles[:len(candles)-1]

	// calculate RSI
	rsi := calculateRSI(&cds)
	// ------------

	// combine all conditions
	isRSIUnder30 := rsi < 30
	acceptedMinStopLoss := 7.0
	// Dùng Max(account.MinStopLoss, acceptedMinStopLoss) để không cho phép MinStopLoss nhỏ hơn 7
	shouldBuy := isRSIUnder30 && (usdtBalance > math.Max(account.MinStopLoss, acceptedMinStopLoss) || usdtBalance > acceptedMinStopLoss)

	fmt.Printf("%s Oldest price: %f\n%s Oldest time: %s\n%s Ask price: %f\n%s Ask value: %f\n%s RSI: %f (<30 - %t)\n%s Should Buy: %t\n%s USDT Balance (>7): %f\n",
		// OLDEST
		prefixLog,
		oldestPrice,
		prefixLog,
		oldestTime,

		// ASK
		prefixLog,
		askPrice,
		prefixLog,
		askValue,

		prefixLog,
		rsi,
		isRSIUnder30,
		prefixLog,
		shouldBuy,
		prefixLog,
		usdtBalance,
	)

	if shouldBuy {
		// check if the 70% of the quantity they want to sell
		// is available based on the amount we bought
		onSellValue := (askValue * askPrice) * 0.7
		quantity := math.Min(usdtBalance/askPrice, onSellValue)
		// ------------
		// calculate the quantity we want to buy
		quantity = utils.FloorTo(quantity, int(account.StepSize))

		fmt.Printf("%s Quantity: %f\n", prefixLog, quantity)
		// ------

		// buy #########################################################
		buyChan := binanceClient.Buy(account, quantity, askPrice, "LIMIT")
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

		purpose := "by _*`down trend`*_"

		if isRSIUnder30 {
			purpose = "by _*`RSI under 30`*_"
		}

		msg := fmt.Sprintf(":%s: :dollar: [BUY] %f (%s) with price *%f* %s - order id: `%d`",
			strings.ToLower(account.Symbol), // emoji
			quantity,
			strings.ToUpper(account.Base),
			askPrice,
			purpose,
			buyResponse.OrderId,
		)

		bodyText := slack.NewTextBlockObject("mrkdwn", msg, false, true)
		bodyBlock := slack.NewSectionBlock(bodyText, nil, nil)

		blocks := []slack.Block{bodyBlock}

		tsChan := slackClient.SendInfo(title, "", blocks...)
		// ------------

		// create stack trade
		ts := <-tsChan

		priceSell := askPrice * (1 + account.Profit) / math.Pow(1-account.Fee, 2)

		now := time.Now()

		quantity4Sell := utils.FloorTo((quantity+baseBalance)*(1-account.Fee), int(account.StepSize))

		// (price_sell <= ? OR price_buy >= ?)
		repositories.NewStackTradeRepository().Create(models.StackTrade{
			Symbol:    account.Symbol,
			Quantity:  quantity4Sell,
			PriceBuy:  askPrice,
			PriceSell: priceSell,
			ThreadID:  ts,
			Status:    "BUY",
			CreatedAt: now,
			UpdatedAt: now,
		})
		// ------------

		titlePriceShouldSell := fmt.Sprintf("💰 Price should sell: %f", priceSell)
		msgPriceShouldSell := fmt.Sprintf("💰 *Price should sell*: `%f`", priceSell)

		bodyTextPriceShouldSell := slack.NewTextBlockObject("mrkdwn", msgPriceShouldSell, false, true)
		bodyBlockPriceShouldSell := slack.NewSectionBlock(bodyTextPriceShouldSell, nil, nil)
		blocksPriceShouldSell := []slack.Block{bodyBlockPriceShouldSell}
		<-slackClient.SendInfo(titlePriceShouldSell, ts, blocksPriceShouldSell...)

		return
	}
	fmt.Println(prefixLog + "Process Buy DONE! =======")
}

func onPostSell(account *models.Account, usdtBalance float64) *models.Account {
	stackTrades := repositories.NewStackTradeRepository().FindSymbolOnSell(account.Symbol)

	totalBalance := usdtBalance

	for _, stackTrade := range stackTrades {
		totalBalance += stackTrade.Quantity * stackTrade.PriceSell
	}

	if totalBalance*0.1 > 30.0 {
		account.MinStopLoss = totalBalance - 30.0
	} else {
		account.MinStopLoss = totalBalance * 0.9
	}

	return account
}

func processSell(t string, account *models.Account, bids *[]binance.Bid, usdtBalance float64) {
	prefixLog := fmt.Sprintf("%s SELL_%s: ", t, account.Symbol)
	fmt.Printf("%s Process Sell =======\n", prefixLog)
	for i := len(*bids) - 1; i >= 0; i-- {
		bid := (*bids)[i]
		bidPrice, _ := strconv.ParseFloat(bid.Price, 64)
		bidValue, _ := strconv.ParseFloat(bid.Quantity, 64)

		b70 := bidValue * 0.7

		stackTrades := repositories.NewStackTradeRepository().FindSymbol4Sell(account.Symbol, "BUY", bidPrice, b70)
		var isStopLoss bool

		if len(stackTrades) > 0 {
			stackTrade := stackTrades[0]
			isStopLoss = stackTrade.PriceBuy*(1-account.StopLoss) >= bidPrice

			purpose := "*sell*"
			if isStopLoss {
				purpose = "`stoploss`"
			}

			sellChan := binanceClient.Sell(account, stackTrade.Quantity, bidPrice, "LIMIT")
			sellResponse := <-sellChan

			quantityEarn, _ := strconv.ParseFloat(sellResponse.CummulativeQuoteQty, 64)

			shouldWithdraw := usdtBalance+quantityEarn > account.MaxWithdraw+account.InitialInvestment
			withdrawQuantity := account.InitialInvestment

			if shouldWithdraw {
				<-binanceClient.Withdraw(account, "USDT", withdrawQuantity)
			}

			stackTrade.Status = "SELL"
			stackTrade.PriceSell = bidPrice
			stackTrade.UpdatedAt = time.Now()
			repositories.NewStackTradeRepository().Update(*stackTrade)

			// quote balance = quantity has bought * current price sell
			quoteBalance := stackTrade.Quantity * bidPrice

			account.BaseBalance = math.Max(0, account.BaseBalance-stackTrade.Quantity)
			account.QuoteBalance = quoteBalance

			account = onPostSell(account, usdtBalance+quantityEarn)

			repositories.NewAccountRepository().Update(*account)

			// log to slack
			title := fmt.Sprintf("💰 Sell %f (%s) with %f\n - Balance: %f (USDT)", stackTrade.Quantity, strings.ToUpper(account.Base), bidPrice, quantityEarn)
			msg := fmt.Sprintf(":%s: :dollar: [SELL] %f (%s) with price *%f* - Balance: %f (USDT)\nby %s",
				strings.ToLower(account.Symbol), // emoji
				stackTrade.Quantity,
				strings.ToUpper(account.Base),
				bidPrice,
				quantityEarn,
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

			<-slackClient.SendInfo(title, stackTrade.ThreadID, blocks...)

			return
		}

		fmt.Printf("%s No stack trade found\n%s Bid price: %f\n%s Bid value: %f\n%s BID 70: %f\n%s STOP LOSS: %t - %f\n\n",
			prefixLog,
			// BID
			prefixLog,
			bidPrice,
			prefixLog,
			bidValue,

			// BID 70
			prefixLog,
			b70,

			// STOP LOSS
			prefixLog,
			isStopLoss,
			bidPrice, // STOP LOSS
		)
	}
	fmt.Println(prefixLog + "Process Sell DONE!=======")
}

func start(symbol string, network string, bids *[]binance.Bid, asks *[]binance.Ask) func() {
	return func() {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			fmt.Println(err)
			return
		}
		t := fmt.Sprintf("[%s]", time.Now().In(loc).Format("2006-01-02 15:04:05"))
		// get account from database
		account := repositories.NewAccountRepository().FindBySymbol(&symbol, &network)
		if account == nil {
			fmt.Printf("%s %s - STOP! Account not found\n", t, symbol)
			return
		}

		if account.IsActived == 0 {
			fmt.Printf("%s %s - STOP! Account is not actived\n", t, symbol)
			return
		}
		// ------------

		// get account info and candlestick data
		accountInfoChan := binanceClient.AccountInfo(account)

		if accountInfoChan == nil {
			fmt.Printf("%s %s - STOP! AccountInfo is not available\n", t, symbol)
			return
		}

		accountInfo := <-accountInfoChan
		// ------------

		// usdtBalance only free USDT balance
		usdtBalance := 0.0
		baseBalance := 0.0
		wgBalance := sync.WaitGroup{}
		wgBalance.Add(2)

		// get usdt balance
		go func() {
			defer wgBalance.Done()
			for _, balance := range accountInfo.Balances {
				if balance.Asset == "USDT" {
					usdtBalance, _ = strconv.ParseFloat(balance.Free, 64)
					break
				}
			}
		}()
		go func() {
			defer wgBalance.Done()
			for _, balance := range accountInfo.Balances {
				if balance.Asset == account.Base {
					num, _ := strconv.ParseFloat(balance.Free, 64)
					baseBalance = num
					break
				}
			}
		}()
		wgBalance.Wait()
		// ------------

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			processBuy(t, account, asks, usdtBalance, baseBalance, loc)
		}()
		go func() {
			defer wg.Done()
			// b := utils.ReverseBids(*bids)
			processSell(t, account, bids, usdtBalance)
		}()
		wg.Wait()
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
	msg := fmt.Sprintf("Error: %v", err)

	bodyText := slack.NewTextBlockObject("mrkdwn", msg, false, true)
	errorBlock := []slack.Block{slack.NewSectionBlock(bodyText, nil, nil)}

	slackClient.SendError("Error", "", errorBlock...)
	time.Sleep(2 * time.Second)
	log.Panic(err)
}

func Run() {
	slackClient = &utils.SlackClient{}
	slackClient.NewSlackClient()

	symbolAndNetworkStr := os.Getenv("SYMBOL_AND_NETWORK")

	if symbolAndNetworkStr == "" {
		str := flag.String("symbol", "", "symbol and network (symbol1:network,symbol2:network)")

		flag.Parse()

		if *str == "" {
			fmt.Println("symbol and network is required")
			return
		}

		symbolAndNetworkStr = *str
	}

	websocketStreamClient := binance.NewWebsocketStreamClient(false)

	symbolsAndNetworks := strings.Split(symbolAndNetworkStr, ",")

	doneChs := make([]chan struct{}, len(symbolsAndNetworks))
	wg := sync.WaitGroup{}

	for i, symbolAndNetwork := range symbolsAndNetworks {
		symbolAndNetwork := strings.Split(symbolAndNetwork, ":")

		if len(symbolAndNetwork) < 2 {
			t := fmt.Sprintf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
			fmt.Printf("%s - STOP! Invalid symbol and network: %s\nSymbol should be like this: symbol1:network,symbol2:network\n", t, symbolAndNetwork)
			return
		}
		symbol := strings.ToUpper(strings.TrimSpace(symbolAndNetwork[0]))
		network := strings.TrimSpace(symbolAndNetwork[1])

		wg.Add(1)
		go func(idx int, sym string, net string) {
			defer wg.Done()
			doneCh, _, err := websocketStreamClient.WsPartialDepthServe100Ms(sym, "20", wsDepthHandler(sym, net), errHandler)
			if err != nil {
				errHandler(err)
				return
			}
			doneChs[idx] = doneCh
		}(i, symbol, network)
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
