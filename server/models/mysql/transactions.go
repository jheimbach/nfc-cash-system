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
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions ORDER BY created`

	return t.loadTransactions(selectStmt)
}

func (t *TransactionModel) GetAllPaged(page, size int) (*models.TransactionPaging, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions ORDER BY created DESC LIMIT ? OFFSET ?`

	transactions, err := t.loadTransactions(selectStmt, size, pageOffset(page, size))

	count, err := countAllIds(t.db, "SELECT COUNT(id) FROM transactions")
	if err != nil {
		return nil, err
	}

	return &models.TransactionPaging{
		CurrentPage:  page,
		MaxPage:      maxPageCount(count, size),
		Transactions: transactions,
	}, nil
}

func (t *TransactionModel) GetAllByAccount(accountId int) ([]*models.Transaction, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE account_id=? ORDER BY created DESC`

	return t.loadTransactions(selectStmt, accountId)
}

func (t *TransactionModel) GetAllByAccountPaged(accountId, page, size int) (*models.TransactionPaging, error) {
	getStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE account_id=? ORDER BY created DESC LIMIT ? OFFSET ?`

	transactions, err := t.loadTransactions(getStmt, accountId, size, pageOffset(page, size))

	count, err := countAllIds(t.db, "SELECT COUNT(id) FROM transactions WHERE account_id=?", accountId)
	if err != nil {
		return nil, err
	}

	return &models.TransactionPaging{
		CurrentPage:  page,
		MaxPage:      maxPageCount(count, size),
		Transactions: transactions,
	}, nil
}

func (t *TransactionModel) loadTransactions(query string, args ...interface{}) ([]*models.Transaction, error) {
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		s := &models.Transaction{Account: &models.Account{}}
		err := rows.Scan(&s.ID, &s.NewSaldo, &s.OldSaldo, &s.Amount, &s.Account.ID, &s.Created)
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
