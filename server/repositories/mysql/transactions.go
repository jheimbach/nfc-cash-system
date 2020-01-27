package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/repositories"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
)

// TransactionRepository provides API for the transactions table
type TransactionRepository struct {
	db       *sql.DB
	accounts repositories.AccountStorager
}

func NewTransactionRepository(db *sql.DB, accounts repositories.AccountStorager) *TransactionRepository {
	return &TransactionRepository{
		db:       db,
		accounts: accounts,
	}
}

// Create inserts new Transaction to database
// with account.saldo and amount, the fields OldSaldo and NewSaldo are calculated
// It will return models.ErrAccountNotFound if account with accountId is not found
func (t *TransactionRepository) Create(ctx context.Context, amount float64, accountId int32) (*api.Transaction, error) {
	// load account
	account, err := t.accounts.Read(ctx, accountId)
	if err != nil {
		return nil, repositories.ErrAccountNotFound
	}

	// calculate saldos
	oldSaldo := account.Saldo
	newSaldo := oldSaldo - amount

	// created time
	now := time.Now()
	nowProto, _ := ptypes.TimestampProto(now)

	// create transaction
	insertStatement := `INSERT INTO transactions (new_saldo, old_saldo, amount, account_id, created) VALUES (?,?,?,?,?)`
	res, err := t.db.ExecContext(ctx, insertStatement, newSaldo, oldSaldo, amount, accountId, now)
	if err != nil {
		// theoretically this error should not be true, we recieve the account in the beginning from accountmodel.
		// but someone could create a transactions and someone else deletes the account
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1452 {
			return nil, repositories.ErrAccountNotFound
		}
		return nil, err
	}

	// update account saldo
	err = t.accounts.UpdateSaldo(ctx, account, newSaldo) // in database
	if err != nil {
		return nil, err
	}
	account.Saldo = newSaldo // in object

	// get id for transaction
	lastId, _ := res.LastInsertId()

	// create transaction object and return it
	return &api.Transaction{
		Id:       int32(lastId),
		OldSaldo: oldSaldo,
		NewSaldo: newSaldo,
		Amount:   amount,
		Created:  nowProto,
		Account:  account,
	}, nil
}

// Read returns Transaction with given id, returns models.ErrNotFound if transaction with id does not exist
func (t *TransactionRepository) Read(ctx context.Context, id int32) (*api.Transaction, error) {
	getSmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions WHERE id=?`

	transaction := &api.Transaction{Account: &api.Account{}}
	var created time.Time

	err := t.db.QueryRowContext(ctx, getSmt, id).Scan(
		&transaction.Id, &transaction.NewSaldo, &transaction.OldSaldo,
		&transaction.Amount, &transaction.Account.Id, &created,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}

	createdProto, err := ptypes.TimestampProto(created)
	if err != nil {
		return nil, err
	}
	transaction.Created = createdProto

	account, err := t.accounts.Read(ctx, transaction.Account.Id)
	if err != nil {
		return nil, err
	}
	transaction.Account = account

	return transaction, nil
}

// GetAll returns all transactions ordered by create date with parameter `order` can be changed (default DESC)
// CAUTION: due to the nature of Transactions, this could be a lot
func (t *TransactionRepository) GetAll(ctx context.Context, accountId int32, order string, limit, offset int32) ([]*api.Transaction, int, error) {
	selectStmt := `SELECT id, new_saldo, old_saldo, amount, account_id, created FROM transactions`

	var args []interface{}
	if accountId > 0 {
		selectStmt = fmt.Sprintf("%s WHERE account_id = ?", selectStmt)
		args = append(args, accountId)
	}

	selectStmt = orderByClause(order, selectStmt)
	if limit > 0 {
		selectStmt = fmt.Sprintf("%s LIMIT ?", selectStmt)
		args = append(args, limit)
		if offset > 0 {
			selectStmt = fmt.Sprintf("%s OFFSET ?", selectStmt)
			args = append(args, offset)
		}
	}

	rows, err := t.db.QueryContext(ctx, selectStmt, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions, err := t.loadTransactions(ctx, rows)

	totalCount := len(transactions)
	if limit > 0 {
		totalCount, err = t.countAll(ctx, accountId)
		if err != nil {
			return nil, 0, err
		}
	}

	return transactions, totalCount, err
}

// DeleteAllByAccount deletes all transactions for given account id
func (t *TransactionRepository) DeleteAllByAccount(ctx context.Context, accountId int32) error {
	delStmt := "DELETE FROM transactions WHERE account_id=?"
	_, err := t.db.ExecContext(ctx, delStmt, accountId)

	return err
}

// orderByClause returns selectStmt with order by created attached.
// If order is ASC or asc returns ORDER BY created ASC, otherwise DESC
func orderByClause(order string, selectStmt string) string {
	switch {
	case strings.ToLower(order) == "asc":
		order = "ASC"
	default:
		order = "DESC"
	}
	selectStmt = fmt.Sprintf("%s ORDER BY created %s", selectStmt, order)
	return selectStmt
}

// countAll counts the account rows in the database and returns a total count
func (t *TransactionRepository) countAll(ctx context.Context, accountId int32) (int, error) {
	countStmt := `SELECT COUNT(id) FROM transactions`
	var countArgs []interface{}

	// if accountId is set, add WHERE clause to count statement
	if accountId > 0 {
		countStmt = fmt.Sprintf("%s WHERE account_id = ?", countStmt)
		countArgs = append(countArgs, accountId)
	}

	var totalCount int
	err := t.db.QueryRowContext(ctx, countStmt, countArgs...).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// loadTransactions will Transactions for given query
func (t *TransactionRepository) loadTransactions(ctx context.Context, rows *sql.Rows) ([]*api.Transaction, error) {

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

	accounts, err := t.accounts.GetAllByIds(ctx, accountIds)
	if err != nil {
		return nil, err
	}

	for _, transaction := range transactions {
		transaction.Account = accounts[transaction.Account.Id]
	}

	return transactions, nil
}
