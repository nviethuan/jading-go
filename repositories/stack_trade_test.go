package repositories

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/utils"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db := utils.GetDBStackTrade()
	db.Exec("DELETE FROM stack_trades")
	return db
}

func TestFindBySymbol(t *testing.T) {
	db := setupTestDB()
	repo := &StackTradeRepository{db: db}

	// Tạo dữ liệu mẫu
	stackTrade1 := models.StackTrade{
		Symbol:    "BTCUSDT",
		Status:    "BUY",
		PriceBuy:  50000,
		Quantity:  0.1,
		PriceSell: 51000,
		StopLoss:  50500, // cập nhật stoploss khác với priceSell
		ThreadID:  "thread1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	stackTrade2 := models.StackTrade{
		Symbol:    "BTCUSDT",
		Status:    "BUY",
		PriceBuy:  50000,
		Quantity:  0.2,
		PriceSell: 52000,
		StopLoss:  51500, // cập nhật stoploss khác với priceSell
		ThreadID:  "thread2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	stackTrade3 := models.StackTrade{
		Symbol:    "ETHUSDT",
		Status:    "SELL",
		PriceBuy:  3000,
		Quantity:  1,
		PriceSell: 3200,
		StopLoss:  3100, // cập nhật stoploss khác với priceSell
		ThreadID:  "thread3",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(stackTrade1)
	repo.Create(stackTrade2)
	repo.Create(stackTrade3)

	// Test tìm kiếm với symbol, status, priceSell, quantity phù hợp
	results := repo.FindBySymbol("BTCUSDT", "BUY", 52000, 0.2, 0.1)

	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Errorf("Lỗi khi chuyển đổi kết quả sang JSON: %v", err)
	} else {
		fmt.Println("Kết quả FindBySymbol (dạng JSON):", string(b))
	}

	if len(results) != 2 {
		t.Errorf("Kỳ vọng tìm thấy 2 stack trade, nhưng chỉ tìm thấy %d", len(results))
	}

	// Test tìm kiếm với điều kiện priceSell nhỏ hơn
	results2 := repo.FindBySymbol("BTCUSDT", "BUY", 51000, 0.2, 0.1)
	b, err = json.MarshalIndent(results2, "", "  ")
	if err != nil {
		t.Errorf("Lỗi khi chuyển đổi kết quả sang JSON: %v", err)
	} else {
		fmt.Println("Kết quả FindBySymbol (dạng JSON):", string(b))
	}
	if len(results2) != 1 {
		t.Errorf("Kỳ vọng tìm thấy 1 stack trade, nhưng tìm thấy %d", len(results2))
	}

	// Test tìm kiếm với stoploss lớn hơn priceSell (case stoploss)
	resultsStopLoss := repo.FindBySymbol("BTCUSDT", "BUY", 50000, 0.2, 0.1)
	b, err = json.MarshalIndent(resultsStopLoss, "", "  ")
	if err != nil {
		t.Errorf("Lỗi khi chuyển đổi kết quả sang JSON: %v", err)
	} else {
		fmt.Println("Kết quả FindBySymbol với stoploss (dạng JSON):", string(b))
	}
	if len(resultsStopLoss) != 1 {
		t.Errorf("Kỳ vọng tìm thấy 1 stack trade với stoploss, nhưng tìm thấy %d", len(resultsStopLoss))
	}

	// Test tìm kiếm với symbol không tồn tại
	results3 := repo.FindBySymbol("BNBUSDT", "BUY", 1000, 1, 0.1)
	b, err = json.MarshalIndent(results3, "", "  ")
	if err != nil {
		t.Errorf("Lỗi khi chuyển đổi kết quả sang JSON: %v", err)
	} else {
		fmt.Println("Kết quả FindBySymbol (dạng JSON):", string(b))
	}
	if len(results3) != 0 {
		t.Errorf("Kỳ vọng không tìm thấy stack trade nào, nhưng lại tìm thấy %d", len(results3))
	}
}
