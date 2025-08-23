package repositories

import (
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
		db: utils.GetDB(),
	}
}

func (s *StackTradeRepository) Create(stackTrade models.StackTrade) models.StackTrade {
	now := time.Now()
	stackTrade.CreatedAt = now
	stackTrade.UpdatedAt = now
	s.db.Create(&stackTrade)
	return stackTrade
}

func (s *StackTradeRepository) Update(stackTrade models.StackTrade) models.StackTrade {
	now := time.Now()
	stackTrade.UpdatedAt = now
	s.db.Save(&stackTrade)
	return stackTrade
}

func (s *StackTradeRepository) FindAll() []models.StackTrade {
	var stackTrades []models.StackTrade
	s.db.Find(&stackTrades)
	return stackTrades
}

func (s *StackTradeRepository) FindBySymbol(symbol string) []models.StackTrade {
	var stackTrades []models.StackTrade
	s.db.Order("price ASC").Where("symbol = ?", symbol).Find(&stackTrades)
	return stackTrades
}