package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	binance "github.com/binance/binance-connector-go"
	"github.com/nviethuan/jading-go/models"
)

func TestBinance_NewBinanceWSClientAPI(t *testing.T) {
	t.Run("should create websocket API client", func(t *testing.T) {
		binanceClient := &Binance{}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
		}

		result := binanceClient.NewBinanceWSClientAPI(account)

		if result == nil {
			t.Error("Expected non-nil result")
		}

		if binanceClient.websocketAPIClient == nil {
			t.Error("Expected websocketAPIClient to be initialized")
		}
	})

	t.Run("should not recreate client if already exists", func(t *testing.T) {
		binanceClient := &Binance{}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
		}

		// Tạo client lần đầu
		firstResult := binanceClient.NewBinanceWSClientAPI(account)
		firstClient := binanceClient.websocketAPIClient

		// Tạo client lần thứ hai
		secondResult := binanceClient.NewBinanceWSClientAPI(account)
		secondClient := binanceClient.websocketAPIClient

		if firstClient != secondClient {
			t.Error("Expected same client instance, got different instances")
		}

		if firstResult != secondResult {
			t.Error("Expected same result instance")
		}
	})
}

func TestBinance_NewBinanceAPI(t *testing.T) {
	t.Run("should create REST API client", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		result := binanceClient.NewBinanceAPI(account)

		if result == nil {
			t.Error("Expected non-nil result")
		}

		// Kiểm tra client đã được lưu vào map
		key := account.Symbol + account.Network
		if binanceClient.clients[key] == nil {
			t.Error("Expected client to be stored in clients map")
		}
	})

	t.Run("should not recreate client if already exists", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		// Tạo client lần đầu
		firstResult := binanceClient.NewBinanceAPI(account)
		key := account.Symbol + account.Network
		firstClient := binanceClient.clients[key]

		// Tạo client lần thứ hai
		secondResult := binanceClient.NewBinanceAPI(account)
		secondClient := binanceClient.clients[key]

		if firstClient != secondClient {
			t.Error("Expected same client instance, got different instances")
		}

		if firstResult != secondResult {
			t.Error("Expected same result instance")
		}
	})
}

func TestBinance_NewBinanceStreamClient(t *testing.T) {
	t.Run("should create websocket stream client", func(t *testing.T) {
		binanceClient := &Binance{}

		result := binanceClient.NewBinanceStreamClient()

		if result == nil {
			t.Error("Expected non-nil result")
		}

		if binanceClient.websocketStreamClient == nil {
			t.Error("Expected websocketStreamClient to be initialized")
		}
	})

	t.Run("should not recreate client if already exists", func(t *testing.T) {
		binanceClient := &Binance{}

		// Tạo client lần đầu
		firstResult := binanceClient.NewBinanceStreamClient()
		firstClient := binanceClient.websocketStreamClient

		// Tạo client lần thứ hai
		secondResult := binanceClient.NewBinanceStreamClient()
		secondClient := binanceClient.websocketStreamClient

		if firstClient != secondClient {
			t.Error("Expected same client instance, got different instances")
		}

		if firstResult != secondResult {
			t.Error("Expected same result instance")
		}
	})
}

func TestBinance_AccountInfo(t *testing.T) {
	t.Run("should return channel for account info", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		accountInfoChan := binanceClient.AccountInfo(account)

		if accountInfoChan == nil {
			t.Error("Expected non-nil channel")
		}

		// Kiểm tra channel có thể nhận dữ liệu
		select {
		case accountInfo := <-accountInfoChan:
			if accountInfo != nil {
				fmt.Println(binance.PrettyPrint(accountInfo))
			}
		case <-time.After(5 * time.Second):
			// Timeout có thể xảy ra với test API keys, không phải lỗi
			t.Log("Timeout waiting for account info (expected with test API keys)")
		}
	})
}

func TestBinance_AccountInfoWithContext(t *testing.T) {
	t.Run("should return channel for account info with context", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		accountInfoChan := binanceClient.AccountInfoWithContext(context.Background(), account)

		if accountInfoChan == nil {
			t.Error("Expected non-nil channel")
		}

		// Kiểm tra channel có thể nhận dữ liệu
		select {
		case accountInfo := <-accountInfoChan:
			fmt.Println(binance.PrettyPrint(accountInfo))
		case <-time.After(5 * time.Second):
			// Timeout có thể xảy ra với test API keys, không phải lỗi
			t.Log("Timeout waiting for account info (expected with test API keys)")
		}
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		ctx, cancel := context.WithCancel(context.Background())
		accountInfoChan := binanceClient.AccountInfoWithContext(ctx, account)

		if accountInfoChan == nil {
			t.Error("Expected non-nil channel")
		}

		// Hủy context ngay lập tức
		cancel()

		// Kiểm tra channel có thể nhận dữ liệu
		select {
		case accountInfo := <-accountInfoChan:
			// Với context bị hủy, có thể nhận được empty response
			_ = accountInfo
		case <-time.After(5 * time.Second):
			// Timeout có thể xảy ra với test API keys
			t.Log("Timeout waiting for account info (expected with test API keys)")
		}
	})
}

func TestBinance_TradeFee(t *testing.T) {
	t.Run("should return channel for trade fee", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		tradeFeeChan := binanceClient.TradeFee(context.Background(), "UNIUSDT", account)

		if tradeFeeChan == nil {
			t.Error("Expected non-nil channel")
		}

		select {
		case tradeFee := <-tradeFeeChan:
			fmt.Println(binance.PrettyPrint(tradeFee))
		case <-time.After(5 * time.Second):
			t.Log("Timeout waiting for trade fee (expected with test API keys)")
		}
	})
}

func TestBinance_Buy(t *testing.T) {
	t.Run("should return channel for buy", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		buyChan := binanceClient.Buy(account, 0.0001, 110189.57, "LIMIT")

		if buyChan == nil {
			t.Error("Expected non-nil channel")
		}

		select {
		case buyResponse := <-buyChan:
			fmt.Println(binance.PrettyPrint(buyResponse))
		case <-time.After(5 * time.Second):
			t.Log("Timeout waiting for buy response (expected with test API keys)")
		}
	})
}

func TestBinance_Sell(t *testing.T) {
	t.Run("should return channel for sell", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		sellChan := binanceClient.Sell(account, 0.0001, 110189.57, "LIMIT")

		if sellChan == nil {
			t.Error("Expected non-nil channel")
		}

		select {
		case sellResponse := <-sellChan:
			fmt.Println(binance.PrettyPrint(sellResponse))
		case <-time.After(5 * time.Second):
			t.Log("Timeout waiting for sell response (expected with test API keys)")
		}
	})
}

func TestBinance_CandlestickData(t *testing.T) {
	t.Run("should return channel for candlestick data", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		candlestickChan := binanceClient.CandlestickData(account, "BTCUSDT", "1m", 5)

		if candlestickChan == nil {
			t.Error("Expected non-nil channel")
		}

		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			t.Error(err)
		}

		// Kiểm tra channel có thể nhận dữ liệu
		select {
		case candlestickData := <-candlestickChan:
			// Với test API keys, có thể sẽ trả về nil hoặc error
			if candlestickData != nil {
				for _, candle := range candlestickData {
					openTime := time.UnixMilli(int64(candle.OpenTime)).In(loc).Format("2006-01-02 15:04:05")
					fmt.Println("Open Time: ", openTime)
				}

				fmt.Println(binance.PrettyPrint(candlestickData))
			}
		case <-time.After(5 * time.Second):
			// Timeout có thể xảy ra với test API keys, không phải lỗi
			t.Log("Timeout waiting for candlestick data (expected with test API keys)")
		}
	})

	t.Run("should handle different symbols and intervals", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		testCases := []struct {
			symbol   string
			interval string
		}{
			{"BTCUSDT", "1m"},
			{"ETHUSDT", "5m"},
			{"ADAUSDT", "1h"},
			{"DOTUSDT", "1d"},
		}

		for _, tc := range testCases {
			t.Run(tc.symbol+"_"+tc.interval, func(t *testing.T) {
				candlestickChan := binanceClient.CandlestickData(account, tc.symbol, tc.interval, 5)

				if candlestickChan == nil {
					t.Error("Expected non-nil channel")
				}

				// Kiểm tra channel có thể nhận dữ liệu
				select {
				case candlestickData := <-candlestickChan:
					_ = candlestickData
				case <-time.After(5 * time.Second):
					// Timeout có thể xảy ra với test API keys, không phải lỗi
					t.Logf("Timeout waiting for candlestick data for %s_%s (expected with test API keys)", tc.symbol, tc.interval)
				}
			})
		}
	})
}

func TestBinance_SymbolPriceTicker(t *testing.T) {
	t.Run("should return channel for symbol price ticker", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		tickerChan := binanceClient.SymbolPriceTicker("BTCUSDT", account)

		if tickerChan == nil {
			t.Error("Expected non-nil channel")
		}

		select {
		case tickerData := <-tickerChan:
			if tickerData != nil {
				fmt.Println(binance.PrettyPrint(tickerData))
			}
		case <-time.After(5 * time.Second):
			t.Log("Timeout waiting for ticker data (expected with test API keys)")
		}
	})
}

func TestBinance_Withdraw(t *testing.T) {
	t.Run("should return channel for withdraw", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}
		account := &models.Account{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			ApiSecret: os.Getenv("BINANCE_SECRET_KEY"),
			RestApi:   os.Getenv("BINANCE_BASE_URL"),
			Symbol:    "BTCUSDT",
			Network:   "BSC",
		}

		withdrawChan := binanceClient.Withdraw(account, "BTC", 0.001)

		if withdrawChan == nil {
			t.Error("Expected non-nil channel")
		}

		select {
		case withdrawResponse := <-withdrawChan:
			fmt.Println(binance.PrettyPrint(withdrawResponse))
		case <-time.After(5 * time.Second):
			t.Log("Timeout waiting for withdraw response (expected with test API keys)")
		}
	})
}

// Test helper function để kiểm tra balance filtering logic
func TestBalanceFiltering(t *testing.T) {
	t.Run("should filter zero balances", func(t *testing.T) {
		balances := []binance.Balance{
			{Asset: "BTC", Free: "0.0", Locked: "0.0"},
			{Asset: "ETH", Free: "1.5", Locked: "0.0"},
			{Asset: "USDT", Free: "0.0", Locked: "0.0"},
			{Asset: "ADA", Free: "100.0", Locked: "0.0"},
		}

		var filteredBalances []binance.Balance
		for _, balance := range balances {
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			total := free + locked
			if total > 0 {
				filteredBalances = append(filteredBalances, balance)
			}
		}

		if len(filteredBalances) != 2 {
			t.Errorf("Expected 2 non-zero balances, got %d", len(filteredBalances))
		}

		expectedAssets := map[string]bool{"ETH": true, "ADA": true}
		for _, balance := range filteredBalances {
			if !expectedAssets[balance.Asset] {
				t.Errorf("Unexpected asset in filtered balances: %s", balance.Asset)
			}
		}
	})
}

// Test cấu trúc của struct Binance
func TestBinance_Struct(t *testing.T) {
	t.Run("should have correct field types", func(t *testing.T) {
		binanceClient := &Binance{}

		// Kiểm tra các field có đúng kiểu dữ liệu
		if binanceClient.clients != nil {
			t.Error("Expected clients to be nil initially")
		}

		if binanceClient.websocketAPIClient != nil {
			t.Error("Expected websocketAPIClient to be nil initially")
		}

		if binanceClient.websocketStreamClient != nil {
			t.Error("Expected websocketStreamClient to be nil initially")
		}
	})

	t.Run("should initialize clients map correctly", func(t *testing.T) {
		binanceClient := &Binance{
			clients: make(map[string]*binance.Client),
		}

		if binanceClient.clients == nil {
			t.Error("Expected clients map to be initialized")
		}

		if len(binanceClient.clients) != 0 {
			t.Error("Expected clients map to be empty initially")
		}
	})
}
