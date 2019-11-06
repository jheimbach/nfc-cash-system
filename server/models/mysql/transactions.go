package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
)

type TransactionModel struct {
	db *sql.DB
}

func (t *TransactionModel) Create(amount, oldSaldo, newSaldo float64, accountId int) error {
	return nil
}

func (t *TransactionModel) Read(id int) (*models.Transaction, error) {
	return nil, nil
}

func (t *TransactionModel) Delete(id int) error {
	return nil
}

func (t *TransactionModel) GetAll() ([]*models.Transaction, error) {
	return nil, nil
}

func (t *TransactionModel) GetAllPaged(page, size int) (*models.TransactionPaging, error) {
	return nil, nil
}

func (t *TransactionModel) GetAllByAccount(accountId int) ([]*models.Transaction, error) {
	return nil, nil
}

func (t *TransactionModel) GetAllByPaged(accountId int) (*models.TransactionPaging, error) {
	return nil, nil
}
