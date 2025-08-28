package utils

import (
	// "path/filepath"

	"github.com/nviethuan/jading-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbAccount *gorm.DB
var dbStackTrade *gorm.DB

func init() {
	var err error
	dbPathAccount := "/app/data/accounts.db"
	dbPathStackTrade := "/app/data/stack_trade.db"
	dbAccount, err = gorm.Open(sqlite.Open(dbPathAccount), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbStackTrade, err = gorm.Open(sqlite.Open(dbPathStackTrade), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbAccount.AutoMigrate(&models.Account{})
	dbStackTrade.AutoMigrate(&models.StackTrade{})
}

func GetDBAccount() *gorm.DB {
	return dbAccount
}

func GetDBStackTrade() *gorm.DB {
	return dbStackTrade
}
