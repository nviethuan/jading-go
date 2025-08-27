package repositories

import (
	"testing"

	"github.com/nviethuan/jading-go/models"
)

func TestAccountRepository_Update(t *testing.T) {
	t.Run("should update account successfully", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account test
		account := &models.Account{
			Symbol:            "TESTUSDT",
			Network:           "Testnet",
			Description:       "Test Account for Update",
			Email:             "test-update@test.com",
			ApiKey:            "test-api-key-update",
			ApiSecret:         "test-secret-key-update",
			RestApi:           "https://test.binance.vision",
			WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
			WsStream:          "wss://stream.testnet.binance.vision",
			Base:              "TEST",
			Quote:             "USDT",
			Fee:               0.001,
			Profit:            1.5,
			IsActived:         1,
			BuyQuantity:       10.0,
			MaxWithdraw:       100,
			InitialInvestment: 100,
			StepSize:          1,
		}

		// Tạo account trước
		createdAccount := repo.Create(account)

		// Cập nhật thông tin
		createdAccount.Description = "Updated Description"
		createdAccount.Profit = 2.0
		createdAccount.BuyQuantity = 15.0
		createdAccount.IsActived = 0

		// Thực hiện update
		updatedAccount := repo.Update(createdAccount)

		// Kiểm tra kết quả
		if updatedAccount.ID == 0 {
			t.Error("Expected non-zero ID after update")
		}

		if updatedAccount.Description != "Updated Description" {
			t.Errorf("Expected description to be 'Updated Description', got '%s'", updatedAccount.Description)
		}

		if updatedAccount.Profit != 2.0 {
			t.Errorf("Expected profit to be 2.0, got %f", updatedAccount.Profit)
		}

		if updatedAccount.BuyQuantity != 15.0 {
			t.Errorf("Expected buy price to be 15.0, got %f", updatedAccount.BuyQuantity)
		}

		if updatedAccount.IsActived != 0 {
			t.Errorf("Expected IsActived to be 0, got %d", updatedAccount.IsActived)
		}

		// Kiểm tra UpdatedAt được cập nhật
		if updatedAccount.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set after update")
		}

		// Cleanup - xóa account test
		repo.Delete(updatedAccount)
	})

	t.Run("should update account with all fields", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account test
		account := &models.Account{
			Symbol:            "TESTUSDT2",
			Network:           "Testnet",
			Description:       "Test Account for Full Update",
			Email:             "test-full-update@test.com",
			ApiKey:            "test-api-key-full-update",
			ApiSecret:         "test-secret-key-full-update",
			RestApi:           "https://test.binance.vision",
			WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
			WsStream:          "wss://stream.testnet.binance.vision",
			Base:              "TEST",
			Quote:             "USDT",
			Fee:               0.001,
			Profit:            1.0,
			IsActived:         1,
			BuyQuantity:       5.0,
			MaxWithdraw:       50,
			InitialInvestment: 50,
			StepSize:          1,
		}

		// Tạo account trước
		createdAccount := repo.Create(account)

		// Cập nhật tất cả các field có thể thay đổi
		createdAccount.Description = "Fully Updated Description"
		createdAccount.Email = "updated-email@test.com"
		createdAccount.ApiKey = "updated-api-key"
		createdAccount.ApiSecret = "updated-secret-key"
		createdAccount.RestApi = "https://updated.binance.vision"
		createdAccount.WsApi = "wss://updated-ws-api.testnet.binance.vision/ws-api/v3"
		createdAccount.WsStream = "wss://updated-stream.testnet.binance.vision"
		createdAccount.Base = "UPDATED"
		createdAccount.Quote = "BTC"
		createdAccount.Fee = 0.002
		createdAccount.Profit = 3.0
		createdAccount.BuyQuantity = 20.0
		createdAccount.MaxWithdraw = 200
		createdAccount.InitialInvestment = 200
		createdAccount.StepSize = 3
		createdAccount.IsActived = 0

		// Thực hiện update
		updatedAccount := repo.Update(createdAccount)

		// Kiểm tra tất cả các field được cập nhật
		if updatedAccount.Description != "Fully Updated Description" {
			t.Errorf("Expected description to be 'Fully Updated Description', got '%s'", updatedAccount.Description)
		}

		if updatedAccount.Email != "updated-email@test.com" {
			t.Errorf("Expected email to be 'updated-email@test.com', got '%s'", updatedAccount.Email)
		}

		if updatedAccount.ApiKey != "updated-api-key" {
			t.Errorf("Expected API key to be 'updated-api-key', got '%s'", updatedAccount.ApiKey)
		}

		if updatedAccount.Base != "UPDATED" {
			t.Errorf("Expected base to be 'UPDATED', got '%s'", updatedAccount.Base)
		}

		if updatedAccount.Quote != "BTC" {
			t.Errorf("Expected quote to be 'BTC', got '%s'", updatedAccount.Quote)
		}

		if updatedAccount.Fee != 0.002 {
			t.Errorf("Expected fee to be 0.002, got %f", updatedAccount.Fee)
		}

		if updatedAccount.Profit != 3.0 {
			t.Errorf("Expected profit to be 3.0, got %f", updatedAccount.Profit)
		}

		if updatedAccount.BuyQuantity != 20.0 {
			t.Errorf("Expected buy price to be 20.0, got %f", updatedAccount.BuyQuantity)
		}

		if updatedAccount.MaxWithdraw != 200 {
			t.Errorf("Expected max withdraw to be 200, got %f", updatedAccount.MaxWithdraw)
		}

		if updatedAccount.InitialInvestment != 200 {
			t.Errorf("Expected initial investment to be 200, got %f", updatedAccount.InitialInvestment)
		}

		if updatedAccount.StepSize != 3 {
			t.Errorf("Expected step size to be 3, got %d", updatedAccount.StepSize)
		}

		if updatedAccount.IsActived != 0 {
			t.Errorf("Expected IsActived to be 0, got %d", updatedAccount.IsActived)
		}

		// Cleanup - xóa account test
		repo.Delete(updatedAccount)
	})

	t.Run("should update account with zero values", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account test
		account := &models.Account{
			Symbol:            "TESTUSDT3",
			Network:           "Testnet",
			Description:       "Test Account for Zero Update",
			Email:             "test-zero-update@test.com",
			ApiKey:            "test-api-key-zero-update",
			ApiSecret:         "test-secret-key-zero-update",
			RestApi:           "https://test.binance.vision",
			WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
			WsStream:          "wss://stream.testnet.binance.vision",
			Base:              "TEST",
			Quote:             "USDT",
			Fee:               0.001,
			Profit:            1.0,
			IsActived:         1,
			BuyQuantity:       10.0,
			MaxWithdraw:       100,
			InitialInvestment: 100,
			StepSize:          2,
		}

		// Tạo account trước
		createdAccount := repo.Create(account)

		// Cập nhật với zero values
		createdAccount.Profit = 0.0
		createdAccount.BuyQuantity = 0.0
		createdAccount.MaxWithdraw = 0.0
		createdAccount.InitialInvestment = 0.0
		createdAccount.StepSize = 0.0
		createdAccount.Fee = 0.0

		// Thực hiện update
		updatedAccount := repo.Update(createdAccount)

		// Kiểm tra zero values được lưu
		if updatedAccount.Profit != 0.0 {
			t.Errorf("Expected profit to be 0.0, got %f", updatedAccount.Profit)
		}

		if updatedAccount.BuyQuantity != 0.0 {
			t.Errorf("Expected buy price to be 0.0, got %f", updatedAccount.BuyQuantity)
		}

		if updatedAccount.MaxWithdraw != 0.0 {
			t.Errorf("Expected max withdraw to be 0.0, got %f", updatedAccount.MaxWithdraw)
		}

		if updatedAccount.InitialInvestment != 0.0 {
			t.Errorf("Expected initial investment to be 0.0, got %f", updatedAccount.InitialInvestment)
		}

		if updatedAccount.StepSize != 0 {
			t.Errorf("Expected step size to be 0.0, got %d", updatedAccount.StepSize)
		}

		if updatedAccount.Fee != 0.0 {
			t.Errorf("Expected fee to be 0.0, got %f", updatedAccount.Fee)
		}

		// Cleanup - xóa account test
		repo.Delete(updatedAccount)
	})
}

func TestAccountRepository_Delete(t *testing.T) {
	t.Run("should delete account successfully", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account test để xóa
		account := &models.Account{
			Symbol:            "TESTUSDT4",
			Network:           "Testnet",
			Description:       "Test Account for Delete",
			Email:             "test-delete@test.com",
			ApiKey:            "test-api-key-delete",
			ApiSecret:         "test-secret-key-delete",
			RestApi:           "https://test.binance.vision",
			WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
			WsStream:          "wss://stream.testnet.binance.vision",
			Base:              "TEST",
			Quote:             "USDT",
			Fee:               0.001,
			Profit:            1.0,
			IsActived:         1,
			BuyQuantity:       10.0,
			MaxWithdraw:       100,
			InitialInvestment: 100,
			StepSize:          2,
		}

		// Tạo account trước
		createdAccount := repo.Create(account)

		// Kiểm tra account đã được tạo
		if createdAccount.ID == 0 {
			t.Error("Expected non-zero ID after creation")
		}

		// Thực hiện xóa
		deletedAccount := repo.Delete(createdAccount)

		// Kiểm tra account đã được xóa
		if deletedAccount.ID == 0 {
			t.Error("Expected non-zero ID in deleted account")
		}

		// Kiểm tra account không còn tồn tại trong database
		foundAccount := repo.FindByID(int(createdAccount.ID))
		if foundAccount.ID != 0 {
			t.Error("Expected account to be deleted from database")
		}
	})

	t.Run("should delete account by symbol and network", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account test với symbol và network cụ thể
		account := &models.Account{
			Symbol:            "TESTUSDT5",
			Network:           "Testnet",
			Description:       "Test Account for Delete by Symbol",
			Email:             "test-delete-symbol@test.com",
			ApiKey:            "test-api-key-delete-symbol",
			ApiSecret:         "test-secret-key-delete-symbol",
			RestApi:           "https://test.binance.vision",
			WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
			WsStream:          "wss://stream.testnet.binance.vision",
			Base:              "TEST",
			Quote:             "USDT",
			Fee:               0.001,
			Profit:            1.0,
			IsActived:         1,
			BuyQuantity:       10.0,
			MaxWithdraw:       100,
			InitialInvestment: 100,
			StepSize:          2,
		}

		// Tạo account trước
		createdAccount := repo.Create(account)

		// Kiểm tra account có thể tìm thấy bằng symbol và network
		foundAccount := repo.FindBySymbol(&createdAccount.Symbol, &createdAccount.Network)
		if foundAccount == nil {
			t.Error("Expected to find account by symbol and network before deletion")
		}

		// Thực hiện xóa
		deletedAccount := repo.Delete(createdAccount)

		// Kiểm tra account đã được xóa
		if deletedAccount.ID == 0 {
			t.Error("Expected non-zero ID in deleted account")
		}

		// Kiểm tra account không còn tìm thấy bằng symbol và network
		foundAccountAfterDelete := repo.FindBySymbol(&createdAccount.Symbol, &createdAccount.Network)
		if foundAccountAfterDelete != nil {
			t.Error("Expected account to be deleted and not found by symbol and network")
		}
	})

	t.Run("should handle delete non-existent account", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo account không tồn tại
		nonExistentAccount := models.Account{
			ID:      999999, // ID không tồn tại
			Symbol:  "NONEXISTENT",
			Network: "Testnet",
		}

		// Thực hiện xóa account không tồn tại
		deletedAccount := repo.Delete(nonExistentAccount)

		// Kiểm tra không có lỗi xảy ra
		if deletedAccount.ID != 999999 {
			t.Errorf("Expected deleted account ID to be 999999, got %d", deletedAccount.ID)
		}
	})

	t.Run("should delete multiple accounts", func(t *testing.T) {
		repo := NewAccountRepository()

		// Tạo nhiều account test
		accounts := []*models.Account{
			{
				Symbol:            "TESTUSDT6",
				Network:           "Testnet",
				Description:       "Test Account 1 for Multiple Delete",
				Email:             "test-multi-delete1@test.com",
				ApiKey:            "test-api-key-multi-delete1",
				ApiSecret:         "test-secret-key-multi-delete1",
				RestApi:           "https://test.binance.vision",
				WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
				WsStream:          "wss://stream.testnet.binance.vision",
				Base:              "TEST",
				Quote:             "USDT",
				Fee:               0.001,
				Profit:            1.0,
				IsActived:         1,
				BuyQuantity:       10.0,
				MaxWithdraw:       100,
				InitialInvestment: 100,
				StepSize:          2,
			},
			{
				Symbol:            "TESTUSDT7",
				Network:           "Testnet",
				Description:       "Test Account 2 for Multiple Delete",
				Email:             "test-multi-delete2@test.com",
				ApiKey:            "test-api-key-multi-delete2",
				ApiSecret:         "test-secret-key-multi-delete2",
				RestApi:           "https://test.binance.vision",
				WsApi:             "wss://ws-api.testnet.binance.vision/ws-api/v3",
				WsStream:          "wss://stream.testnet.binance.vision",
				Base:              "TEST",
				Quote:             "USDT",
				Fee:               0.001,
				Profit:            1.0,
				IsActived:         1,
				BuyQuantity:       10.0,
				MaxWithdraw:       100,
				InitialInvestment: 100,
				StepSize:          2,
			},
		}

		// Tạo các account
		createdAccounts := make([]models.Account, len(accounts))
		for i, account := range accounts {
			createdAccounts[i] = repo.Create(account)
		}

		// Kiểm tra các account đã được tạo
		for i, createdAccount := range createdAccounts {
			if createdAccount.ID == 0 {
				t.Errorf("Expected non-zero ID for account %d after creation", i)
			}
		}

		// Xóa tất cả các account
		for i, createdAccount := range createdAccounts {
			deletedAccount := repo.Delete(createdAccount)

			if deletedAccount.ID == 0 {
				t.Errorf("Expected non-zero ID for deleted account %d", i)
			}

			// Kiểm tra account đã được xóa
			foundAccount := repo.FindByID(int(createdAccount.ID))
			if foundAccount.ID != 0 {
				t.Errorf("Expected account %d to be deleted from database", i)
			}
		}
	})
}
