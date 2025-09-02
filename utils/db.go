package utils

import (
	"os"

	"github.com/nviethuan/jading-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbAccount *gorm.DB
var dbStackTrade *gorm.DB

func init() {
	var err error
	home, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home directory")
	}
	dbPathAccount := home + "/app/data/accounts.db"
	dbPathStackTrade := home + "/app/data/stack_trade.db"

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
