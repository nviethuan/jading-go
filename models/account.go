package models

import "time"

type Account struct {
	ID                uint   `gorm:"primaryKey"`
	Symbol            string `gorm:"uniqueIndex:idx_symbol_network"`
	Network           string `gorm:"uniqueIndex:idx_symbol_network"`
	Description       string
	Email             string `gorm:"uniqueIndex:idx_email_name"`
	ApiKey            string `gorm:"uniqueIndex:idx_email_name"`
	ApiSecret         string `gorm:"uniqueIndex:idx_email_name"`
	RestApi           string
	WsApi             string
	WsStream          string
	Base              string
	Quote             string
	Fee               float64
	Profit            float64
	BuyPrice          float64
	IsActived         int8 `gorm:"default:0"`
	MaxWithdraw       float64
	InitialInvestment float64
	StepSize          float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
