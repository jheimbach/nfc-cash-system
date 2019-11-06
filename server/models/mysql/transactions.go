package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

type TransactionModel struct {
	db *sql.DB
}

func (t *TransactionModel) Create(amount, oldSaldo, newSaldo float64, accountId int) error {

	insertStatement := `INSERT INTO transactions (new_saldo, old_saldo, amount, account_id, created) VALUES (?,?,?,?,UTC_TIMESTAMP)`
	_, err := t.db.Exec(insertStatement, newSaldo, oldSaldo, amount, accountId)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1452 {
			return models.ErrAccountNotFound
		}
		return err
	}

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
