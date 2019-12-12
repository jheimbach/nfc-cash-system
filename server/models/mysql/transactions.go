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
	db       *sql.DB
	accounts models.AccountStorager
}

func NewTransactionModel(db *sql.DB, accounts models.AccountStorager) *TransactionModel {
	return &TransactionModel{
		db:       db,
		accounts: accounts,
	}
}

// Create inserts new Transaction to database will return models.ErrAccountNotFound if account is not associated with account
func (t *TransactionModel) Create(amount, oldSaldo, newSaldo float64, accountId int32) (*api.Transaction, error) {
	account, err := t.accounts.Read(accountId)
	if err != nil {
		return nil, models.ErrAccountNotFound
	}

	now := time.Now()
	nowProto, _ := ptypes.TimestampProto(now)
	insertStatement := `INSERT INTO transactions (new_saldo, old_saldo, amount, account_id, created) VALUES (?,?,?,?,?)`
	res, err := t.db.Exec(insertStatement, newSaldo, oldSaldo, amount, accountId, now)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1452 {
			return nil, models.ErrAccountNotFound
		}
		return nil, err
	}

	lastId, _ := res.LastInsertId()
	transaction := &api.Transaction{
		Id:       int32(lastId),
		OldSaldo: oldSaldo,
		NewSaldo: newSaldo,
		Amount:   amount,
		Created:  nowProto,
		Account:  account,
	}

	return transaction, nil
}

// Read returns Transaction with given id, returns models.ErrNotFound if transaction with id does not exist
func (t *TransactionModel) Read(id int32) (*api.Transaction, error) {
	getSmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE id=?`

	transaction := &api.Transaction{Account: &api.Account{}}
	var created time.Time

	err := t.db.QueryRow(getSmt, id).Scan(
		&transaction.Id, &transaction.NewSaldo, &transaction.OldSaldo,
		&transaction.Amount, &transaction.Account.Id, &created,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	createdProto, _ := ptypes.TimestampProto(created)
	transaction.Created = createdProto

	return transaction, nil
}

// GetAll returns all transactions
// CAUTION: due to the nature of Transactions, this could be a lot
func (t *TransactionModel) GetAll() ([]*api.Transaction, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions ORDER BY created`

	return t.loadTransactions(selectStmt)
}

// GetAllByAccount returns transactions for given account id
func (t *TransactionModel) GetAllByAccount(accountId int32) ([]*api.Transaction, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE account_id=? ORDER BY created DESC`

	return t.loadTransactions(selectStmt, accountId)
}

// loadTransactions will Transactions for given query
func (t *TransactionModel) loadTransactions(query string, args ...interface{}) ([]*api.Transaction, error) {
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*api.Transaction
	var accountIdsLookup = make(map[int32]bool)
	var accountIds []int32

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

		if _, ok := accountIdsLookup[s.Account.Id]; !ok {
			accountIdsLookup[s.Account.Id] = true
			accountIds = append(accountIds, s.Account.Id)
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if len(transactions) == 0 {
		return nil, nil
	}

	accounts, err := t.accounts.GetAllByIds(accountIds)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Account = accounts[transaction.Account.Id]
	}

	return transactions, nil
}
