package repositories

import (
	"github.com/nviethuan/jading-go/models"
	"github.com/nviethuan/jading-go/utils"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		db: utils.GetDB(),
	}
}

func (t *TransactionRepository) Create(transaction models.Transaction) models.Transaction {
	t.db.Create(&transaction)
	return transaction
}

func (t *TransactionRepository) Update(transaction models.Transaction) models.Transaction {
	t.db.Save(&transaction)
	return transaction
}

func (t *TransactionRepository) FindAll() []models.Transaction {
	var transactions []models.Transaction
	t.db.Find(&transactions)
	return transactions
}

func (t *TransactionRepository) FindByID(id int) models.Transaction {
	var transaction models.Transaction
	t.db.Find(&transaction, "id = ?", id)
	return transaction
}

func (t *TransactionRepository) FindBySymbol(symbol string) []models.Transaction {
	var transactions []models.Transaction
	t.db.Find(&transactions, "symbol = ?", symbol)
	return transactions
}
