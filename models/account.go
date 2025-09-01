package models

import "time"

type Account struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Symbol       string `gorm:"uniqueIndex:idx_symbol_network" json:"symbol"`
	Network      string `gorm:"uniqueIndex:idx_symbol_network" json:"network"`
	Description  string `json:"description"`
	Email        string `gorm:"uniqueIndex:idx_email_name" json:"email"`
	ApiKey       string `gorm:"uniqueIndex:idx_email_name" json:"api_key"`
	ApiSecret    string `gorm:"uniqueIndex:idx_email_name" json:"api_secret"`
	Interval     string `json:"interval"`
	RestApi      string `json:"rest_api"`
	WsApi        string `json:"ws_api"`
	WsStream     string `json:"ws_stream"`
	Base         string `json:"base"`
	Quote        string `json:"quote"`
	BaseBalance  float64 `json:"base_balance"`
	QuoteBalance float64 `json:"quote_balance"`
	// Exchange fee: 0.001 = 0.1%
	Fee float64 `json:"fee"`
	// Profit: 0.01 = 1%
	Profit float64 `json:"profit"`
	// Config stop loss and buy quantity: 0.08 = 8% of usdt balance
	StopLoss float64 `json:"stop_loss"`
	// Percentage to check downtrend: 0.015 = 1.5%
	Threshold         float64 `gorm:"default:0.015" json:"threshold"`
	// This is min number of usdt balance can be used to buy, default: 8 (usdt)
	MinStopLoss       float64 `gorm:"default:8.0" json:"min_stop_loss"`
	BuyQuantity       float64
	IsActived         int8 `gorm:"default:0" json:"is_actived"`
	MaxWithdraw       float64 `gorm:"default:5000.0" json:"max_withdraw"`
	InitialInvestment float64 `json:"initial_investment"`
	// This value stores the number of decimal places to round when trading
	// Example: value 0 <=> 1 (no rounding)
	// Example: value 1 <=> 0.1 (round to 1 decimal place)
	// Example: value 2 <=> 0.01 (round to 2 decimal places)
	StepSize          int8 `json:"step_size"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
