package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/repositories"
	"github.com/go-sql-driver/mysql"
)

const accountFields = "id, name, description, saldo, group_id, nfc_chip_uid"

// AccountRepository provides API for the accounts table
type AccountRepository struct {
	db     *sql.DB
	groups repositories.GroupStorager
}

func NewAccountRepository(db *sql.DB, model repositories.GroupStorager) *AccountRepository {
	return &AccountRepository{
		db:     db,
		groups: model,
	}
}

// Create inserts new account it returns error models.ErrGroupNotFound if the groupId is not associated with a group
// it returns models.ErrDuplicateNfcChipId if the provided nfcchipid is already in the database present
func (a *AccountRepository) Create(ctx context.Context, name, description string, startSaldo float64, groupId int32, nfcChipId string) (*api.Account, error) {
	nullDescription := createNullableString(description)

	group, err := a.groups.Read(ctx, groupId)
	if err != nil {
		if err == repositories.ErrNotFound {
			return nil, repositories.ErrGroupNotFound
		}
		return nil, err
	}

	createStmt := `INSERT INTO accounts (name, description, saldo, group_id, nfc_chip_uid) VALUES (?,?,?,?,?)`

	res, err := a.db.ExecContext(ctx, createStmt, name, nullDescription, startSaldo, group.Id, nfcChipId)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				return nil, repositories.ErrDuplicateNfcChipId
			}
		}
		return nil, err
	}

	// mysql implementation of sql.Result returns no error on LastInsertId, so we can ignore it
	lastId, _ := res.LastInsertId()

	return &api.Account{
		Id:          int32(lastId),
		Name:        name,
		Description: description,
		Saldo:       startSaldo,
		NfcChipId:   nfcChipId,
		Group:       group,
	}, nil
}

// Read returns account struct for given id
func (a *AccountRepository) Read(ctx context.Context, id int32) (*api.Account, error) {
	readStmt := `SELECT ` + accountFields + ` FROM accounts WHERE id=?`

	m := &api.Account{}
	var groupId int32
	row := a.db.QueryRowContext(ctx, readStmt, id)
	var nullDesc sql.NullString
	err := row.Scan(&m.Id, &m.Name, &nullDesc, &m.Saldo, &groupId, &m.NfcChipId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	m.Description = decodeNullableString(nullDesc)

	group, err := a.groups.Read(ctx, groupId)
	if err != nil {
		return nil, err
	}
	m.Group = group

	return m, nil
}

// Update saves the (changed) model in the database will return models.ErrGroupNotFound if group id is not associated with a group
func (a *AccountRepository) Update(ctx context.Context, m *api.Account) (*api.Account, error) {
	acc, err := a.Read(ctx, m.Id)
	if err != nil {
		return nil, err
	}

	if m.Saldo != 0 && m.Saldo != acc.Saldo {
		return nil, repositories.ErrUpdateSaldo
	}

	g, err := a.groups.Read(ctx, m.Group.Id)
	if err != nil {
		return nil, repositories.ErrGroupNotFound
	}

	updateStmt := `UPDATE accounts SET name=?, description=?, group_id=?, nfc_chip_uid=? WHERE id=?`

	_, err = a.db.ExecContext(ctx, updateStmt, m.Name, m.Description, m.Group.Id, m.NfcChipId, m.Id)

	if err != nil {
		return nil, err
	}

	acc.Name = m.Name
	acc.Description = m.Description
	acc.NfcChipId = m.NfcChipId
	acc.Group = g

	return acc, nil
}

// Delete deletes a account
func (a *AccountRepository) Delete(ctx context.Context, id int32) error {

	deleteStmt := `DELETE FROM accounts WHERE id=?`

	_, err := a.db.ExecContext(ctx, deleteStmt, id)

	return err
}

// UpdateSaldo provides update method for the saldo field
func (a *AccountRepository) UpdateSaldo(ctx context.Context, m *api.Account, newSaldo float64) error {
	_, err := a.db.ExecContext(ctx, `UPDATE accounts SET saldo=? WHERE id=?`, newSaldo, m.Id)

	return err
}

// GetAll returns slice with all accounts in the database
func (a *AccountRepository) GetAll(ctx context.Context, groupId int32, limit int32, offset int32) ([]*api.Account, int, error) {
	// default select statement
	stmt := `SELECT ` + accountFields + ` FROM accounts`

	// want slice for query we have max 3 entries (groupid, limit, offset), so we can set the capacity
	args := make([]interface{}, 0, 3)

	// if groupId is set (and not the default value of zero)
	// add WHERE clause to select query
	if groupId > 0 {
		stmt = fmt.Sprintf("%s WHERE group_id = ?", stmt)
		args = append(args, groupId)
	}
	// if limit is set
	// add LIMIT clause to select query
	if limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT ?", stmt)
		args = append(args, limit)

		// if limit and offset is set
		// add OFFSET clause to select query
		if offset > 0 {
			stmt = fmt.Sprintf("%s OFFSET ?", stmt)
			args = append(args, offset)
		}
	}

	// get rows from database
	rows, err := a.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// scan rows to account objects
	accounts, err := a.scanRowsToAccounts(ctx, rows)
	if err != nil {
		return nil, 0, err
	}

	// set totalCount to length of accounts slice
	totalCount := len(accounts)

	// if limit is set, ask the database for the total amount
	// we don't have to ask the database if no limit is set, because all account (in this group) are in the account slice
	if limit > 0 {
		totalCount, err = a.countAll(ctx, groupId)
		if err != nil {
			return nil, 0, err
		}
	}

	return accounts, totalCount, nil
}

// countAll counts the account rows in the database and returns a total count
func (a *AccountRepository) countAll(ctx context.Context, groupId int32) (int, error) {
	countStmt := `SELECT COUNT(id) FROM accounts`
	var countArgs []interface{}

	// if groupId is set, add WHERE clause to count statement
	if groupId > 0 {
		countStmt = fmt.Sprintf("%s WHERE group_id = ?", countStmt)
		countArgs = append(countArgs, groupId)
	}

	var totalCount int
	err := a.db.QueryRowContext(ctx, countStmt, countArgs...).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

// GetAllByIds returns map of accounts, is used by to complete objects that are dependent on accounts
func (a *AccountRepository) GetAllByIds(ctx context.Context, ids []int32) (map[int32]*api.Account, error) {
	m := make(map[int32]*api.Account, len(ids))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	stmt := `SELECT ` + accountFields + ` FROM accounts WHERE id IN (?` + strings.Repeat(",?", len(ids)-1) + `)`
	rows, err := a.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts, err := a.scanRowsToAccounts(ctx, rows)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		m[account.Id] = account
	}

	return m, nil
}

// scanRowsToAccounts returns slice of Accounts from given sql.Rows
func (a *AccountRepository) scanRowsToAccounts(ctx context.Context, rows *sql.Rows) ([]*api.Account, error) {
	var accounts []*api.Account
	var groupIds []int32

	for rows.Next() {
		s := &api.Account{Group: &api.Group{}}

		var nullDesc sql.NullString

		err := rows.Scan(&s.Id, &s.Name, &nullDesc, &s.Saldo, &s.Group.Id, &s.NfcChipId)
		if err != nil {
			return nil, err
		}

		s.Description = decodeNullableString(nullDesc)

		groupIds = append(groupIds, s.Group.Id)
		accounts = append(accounts, s)
	}

	groups, err := a.groups.GetAllByIds(ctx, groupIds)

	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		account.Group = groups[account.Group.Id]
	}
	return accounts, nil
}
