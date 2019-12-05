package mysql

import (
	"database/sql"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
)

// TransactionModel provides API for the transactions table
type TransactionModel struct {
	db *sql.DB
}

// Create inserts new Transaction to database will return models.ErrAccountNotFound if accountId is not associated with account
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

// GetAll will returns slice with all transactions
// CAUTION: due to the nature of Transactions, this could be a lot
func (t *TransactionModel) GetAll() ([]*api.Transaction, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions ORDER BY created`

	return t.loadTransactions(selectStmt)
}

// GetAllByAccount returns slice with transactions for given account id
func (t *TransactionModel) GetAllByAccount(accountId int) ([]*api.Transaction, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE account_id=? ORDER BY created DESC`

	return t.loadTransactions(selectStmt, accountId)
}

// loadTransactions will return slice of Transactions for given query
func (t *TransactionModel) loadTransactions(query string, args ...interface{}) ([]*api.Transaction, error) {
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*api.Transaction

	for rows.Next() {
		s := &api.Transaction{Account: &api.Account{}}
		var t time.Time

		err := rows.Scan(&s.Id, &s.NewSaldo, &s.OldSaldo, &s.Amount, &s.Account.Id, &t)
		if err != nil {
			return nil, err
		}

		s.Created, err = ptypes.TimestampProto(t)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return transactions, nil
}
