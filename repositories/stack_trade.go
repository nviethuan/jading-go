package repositories

import (
	"os/exec"
	"time"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/utils"
	"gorm.io/gorm"
)

type StackTradeRepository struct {
	db *gorm.DB
}

func NewStackTradeRepository() *StackTradeRepository {
	return &StackTradeRepository{
		db: utils.GetDBStackTrade(),
	}
}

func (s *StackTradeRepository) Create(stackTrade models.StackTrade) models.StackTrade {
	now := time.Now()
	stackTrade.CreatedAt = now
	stackTrade.UpdatedAt = now
	s.db.Create(&stackTrade)
	exec.Command("ss3", "-key", "stack_trade.db").Start()
	return stackTrade
}

func (s *StackTradeRepository) Update(stackTrade models.StackTrade) models.StackTrade {
	now := time.Now()
	stackTrade.UpdatedAt = now
	s.db.Save(&stackTrade)
	exec.Command("ss3", "-key", "stack_trade.db").Start()

	return stackTrade
}

func (s *StackTradeRepository) FindAll() []models.StackTrade {
	var stackTrades []models.StackTrade
	s.db.Find(&stackTrades)
	return stackTrades
}

func (s *StackTradeRepository) FindBySymbol(symbol string, status string, priceSell float64, quantity float64, stopLoss float64, currentPrice float64) []*models.StackTrade {
	var stackTrades []*models.StackTrade
	s.db.Order("price_buy DESC").Where(
		"symbol = ? AND status = ? AND quantity <= ? AND price_sell <= ?",
		symbol, status, quantity, priceSell, stopLoss, currentPrice,
	).Limit(1).Find(&stackTrades)
	return stackTrades
}
