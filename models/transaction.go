package models

import (
	"time"
)

type Transaction struct {
	ID        uint           `gorm:"primaryKey"`
	Symbol    string         `gorm:"index:tran_symbol_for_search_2025_08_23;not null"`
	Side      string         `gorm:"not null"`
	Quantity  string         `gorm:"not null"`
	Price     string         `gorm:"not null"`
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
