// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	binance "github.com/binance/binance-connector-go"
)

func main() {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_SECRET_KEY")
	restApi := os.Getenv("BINANCE_BASE_URL")
	// if api doesn't have permission to get candlestick data,
	// please check whitelist IP from binance of this API key:
	client := binance.NewClient(apiKey, apiSecret, restApi)
	symbol := "UNIUSDT"
	interval := "1m"
	limit := 15
	response, err := client.NewKlinesService().Symbol(symbol).Interval(interval).Limit(limit).Do(context.Background())

	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for i, candle := range response {
		fmt.Print(i, " - close time: ", time.UnixMilli(int64(candle.CloseTime)).Format("2006-01-02 15:04:05"), " - ")
		fmt.Println("close: ", candle.Close)
	}

	newCandles := response[:len(response)-1]

	fmt.Println("================================================")

	for i, candle := range newCandles {
		fmt.Print(i, " - close time: ", time.UnixMilli(int64(candle.CloseTime)).Format("2006-01-02 15:04:05"), " - ")
		fmt.Println("close: ", candle.Close)
	}
}
