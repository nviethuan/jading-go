package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/repositories"
)

func main() {

	accounts, err := os.ReadFile("accounts.json")
	if err != nil {
		log.Fatal(err)
	}

	var accountData []models.Account
	json.Unmarshal(accounts, &accountData)

	for _, account := range accountData {
		repositories.NewAccountRepository().Create(&account)
	}
}
