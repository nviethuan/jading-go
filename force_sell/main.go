package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nviethuan/jading-go/repositories"
	"github.com/nviethuan/jading-go/utils"
	"github.com/slack-go/slack"
)

var binanceClient *utils.Binance = &utils.Binance{}
var slackClient *utils.SlackClient = &utils.SlackClient{}

// cp to ~/app
// run with command: fsell -symbol=BTCUSDT
// WARNING: Don't forget ~/app/.env file
func main() {
	slackClient.NewSlackClient()

	symbolInput := flag.String("symbol", "", "symbol")
	flag.Parse()

	if *symbolInput == "" {
		fmt.Println("symbol is required")
		return
	}

	symbol := strings.ToUpper(strings.TrimSpace(*symbolInput))
	network := "mainnet"

	account := repositories.NewAccountRepository().FindBySymbol(&symbol, &network)
	if account == nil {
		fmt.Println("account not found")
		return
	}

	stackTrades := repositories.NewStackTradeRepository().FindSymbol4Sell(account.Symbol, "BUY", 100000000, 10000000000)

	for _, stackTrade := range stackTrades {
		sellChan := binanceClient.SellMarket(account, stackTrade.Quantity)
		sellResponse := <-sellChan
		quantityEarn, _ := strconv.ParseFloat(sellResponse.CummulativeQuoteQty, 64)

		stackTrade.Status = "SELL"

		priceSell, _ := strconv.ParseFloat(sellResponse.Price, 64)

		stackTrade.PriceSell = priceSell
		stackTrade.UpdatedAt = time.Now()
		repositories.NewStackTradeRepository().Update(*stackTrade)

		title := fmt.Sprintf("💰 Sell %f (%s) with %f\n - Balance: %f (USDT)", stackTrade.Quantity, strings.ToUpper(account.Base), priceSell, quantityEarn)
		msg := fmt.Sprintf(":%s: :dollar: [SELL] %f (%s) with price *%f* - Balance: %f (USDT)\nby %s",
			strings.ToLower(account.Symbol), // emoji
			stackTrade.Quantity,
			strings.ToUpper(account.Base),
			priceSell,
			quantityEarn,
			"force sell from CLI",
		)

		bodyText := slack.NewTextBlockObject("mrkdwn", msg, false, true)
		bodyBlock1 := slack.NewSectionBlock(bodyText, nil, nil)
		blocks := []slack.Block{bodyBlock1}
		<-slackClient.SendInfo(title, stackTrade.ThreadID, blocks...)
		time.Sleep(1 * time.Second)
	}
}
