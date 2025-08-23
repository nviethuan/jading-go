package utils

import (
	"os"

	binance "github.com/binance/binance-connector-go"
)

var (
	binanceWSClientAPI *binance.WebsocketAPIClient
	binanceAPI         *binance.Client
	// serverTimeOffset   int64 // Offset giữa local time và server time
)

func NewBinanceWSClientAPI() *binance.WebsocketAPIClient {
	if binanceWSClientAPI == nil {
		binanceWSClientAPI = binance.NewWebsocketAPIClient(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"), os.Getenv("BINANCE_WS_BASE_URL"))
	}
	return binanceWSClientAPI
}

func NewBinanceAPI() *binance.Client {
	if binanceAPI == nil {
		binanceAPI = binance.NewClient(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"), os.Getenv("BINANCE_BASE_URL"))
	}
	return binanceAPI
}

// // GetServerTime lấy server timestamp từ Binance
// func GetServerTime() (int64, error) {
// 	client := NewBinanceAPI()
// 	serverTime, err := client.NewServerTimeService().Do(context.Background())
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(serverTime.ServerTime), nil
// }

// // SyncServerTime đồng bộ thời gian với server Binance
// func SyncServerTime() error {
// 	serverTime, err := GetServerTime()
// 	if err != nil {
// 		return err
// 	}

// 	localTime := time.Now().UnixMilli()
// 	serverTimeOffset = serverTime - localTime

// 	return nil
// }

// // GetAdjustedTimestamp trả về timestamp đã được điều chỉnh theo server time
// func GetAdjustedTimestamp() int64 {
// 	return time.Now().UnixMilli() + serverTimeOffset
// }
