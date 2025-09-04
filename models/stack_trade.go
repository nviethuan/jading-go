package models

import (
	"time"
)

type StackTrade struct {
	ID        uint    `gorm:"primaryKey"`
	Symbol    string  `gorm:"index:symbol_for_search_2025_08_23;not null"`
	// string enum: BUY, SELL
	Status    string  `gorm:"index:symbol_for_search_2025_08_23;not null"`
	// If status is BUY, price sell is the price we want to sell.
	// If status is SELL, price sell is the price we sold.
	PriceSell float64 `gorm:"index:symbol_for_search_2025_08_23;not null;"`
	PriceBuy  float64 `gorm:"index:symbol_for_search_2025_08_23;not null"`
	Quantity  float64 `gorm:"index:symbol_for_search_2025_08_23;not null"`
	ThreadID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
