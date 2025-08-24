package repositories

import (
	"time"

	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/utils"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		db: utils.GetDB(),
	}
}

func (a *AccountRepository) FindAll() []models.Account {
	var accounts []models.Account
	a.db.Find(&accounts)
	return accounts
}

func (a *AccountRepository) FindByID(id int) models.Account {
	var account models.Account
	a.db.Find(&account, id)
	return account
}

func (a *AccountRepository) FindBySymbol(symbol string, network string) *models.Account {
	var account models.Account
	err := a.db.Where("symbol = ? AND network = ?", symbol, network).First(&account).Error
	if err != nil {
		return nil
	}
	return &account
}

func (a *AccountRepository) Create(account *models.Account) models.Account {
	now := time.Now()

	account.CreatedAt = now
	account.UpdatedAt = now

	a.db.Create(&account)
	return *account
}

func (a *AccountRepository) Update(account models.Account) models.Account {
	now := time.Now()
	account.UpdatedAt = now

	a.db.Save(&account)
	return account
}

func (a *AccountRepository) Delete(account models.Account) models.Account {
	a.db.Delete(&account)
	return account
}
