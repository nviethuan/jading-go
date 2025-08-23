package utils

import (
	"github.com/nviethuan/jading-go/models"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(".data/mine.db"), &gorm.Config{

	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Account{}, &models.Transaction{}, &models.StackTrade{})
}

func GetDB() *gorm.DB {
	return db
}
