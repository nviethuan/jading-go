package models

import "time"

type Account struct {
	ID           uint   `gorm:"primaryKey"`
	Symbol       string `gorm:"uniqueIndex:idx_symbol_network"`
	Network      string `gorm:"uniqueIndex:idx_symbol_network"`
	Description  string
	Email        string `gorm:"uniqueIndex:idx_email_name"`
	ApiKey       string `gorm:"uniqueIndex:idx_email_name"`
	ApiSecret    string `gorm:"uniqueIndex:idx_email_name"`
	Interval     string
	RestApi      string
	WsApi        string
	WsStream     string
	Base         string
	Quote        string
	BaseBalance  float64
	QuoteBalance float64
	// Exchange fee: 0.001 = 0.1%
	Fee float64
	// Profit: 0.01 = 1%
	Profit float64
	// Config stop loss and buy quantity: 0.08 = 8% of usdt balance
	StopLoss float64
	// Percentage to check downtrend: 0.015 = 1.5%
	Threshold         float64 `gorm:"default:0.015"`
	// This is min number of usdt balance can be used to buy, default: 8 (usdt)
	MinStopLoss       float64 `gorm:"default:8.0"`
	BuyQuantity       float64
	IsActived         int8 `gorm:"default:0"`
	MaxWithdraw       float64 `gorm:"default:5000.0"`
	InitialInvestment float64
	// This value stores the number of decimal places to round when trading
	// Example: value 0 <=> 1 (no rounding)
	// Example: value 1 <=> 0.1 (round to 1 decimal place)
	// Example: value 2 <=> 0.01 (round to 2 decimal places)
	StepSize          int8
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
