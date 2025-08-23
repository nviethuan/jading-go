package models

import (
	"time"
)

type StackTrade struct {
	ID        uint           `gorm:"primaryKey"`
	Symbol    string         `gorm:"index:symbol_for_search_2025_08_23;not null"`
	Quantity  float64         `gorm:"not null"`
	Price     float64         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
