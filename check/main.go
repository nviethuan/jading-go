package main

// import (
// 	"context"
// 	"log"
// 	"os"

// 	"github.com/nviethuan/jading-go/utils"

// 	binance "github.com/binance/binance-connector-go"
// )

// func getExchangeInfo() {
// 	client := utils.NewBinanceAPI(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"), os.Getenv("BINANCE_BASE_URL"))
// 	exchangeInfo, err := client.NewExchangeInfoService().Symbol("BTCUSDT").Do(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	data := []byte(binance.PrettyPrint(exchangeInfo))

// 	err = os.WriteFile("exchange_info.json", data, 0644)
// 	if err != nil {
// 		log.Fatal("Không thể ghi file:", err)
// 	}

// 	log.Println(binance.PrettyPrint(exchangeInfo))
// }

// func main() {
// 	getExchangeInfo()
// }
