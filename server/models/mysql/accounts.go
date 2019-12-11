package mysql

import (
	"database/sql"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

// AccountModel provides API for the accounts table
type AccountModel struct {
	db     *sql.DB
	groups models.GroupStorager
}

func NewAccountModel(db *sql.DB) *AccountModel {
	return &AccountModel{db: db}
}

// Create inserts new account it returns error models.ErrGroupNotFound if the groupId is not associated with a group
// it returns models.ErrDuplicateNfcChipId if the provided nfcchipid is already in the database present
func (a *AccountModel) Create(name, description string, startSaldo float64, groupId int32, nfcChipId string) (*api.Account, error) {
	nullDescription := createNullableString(description)

	group, err := a.groups.Read(groupId)
	if err != nil {
		return nil, err
	}

	createStmt := `INSERT INTO accounts (name, description, saldo, group_id, nfc_chip_uid) VALUES (?,?,?,?,?)`

	res, err := a.db.Exec(createStmt, name, nullDescription, startSaldo, group.Id, nfcChipId)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1452 {
				return nil, models.ErrGroupNotFound
			}
			if err.Number == 1062 {
				return nil, models.ErrDuplicateNfcChipId
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
func (a *AccountModel) Read(id int32) (*api.Account, error) {
	readStmt := `SELECT id, name, description, saldo, group_id, nfc_chip_uid FROM accounts WHERE id=?`

	m := &api.Account{}
	var groupId int32
	row := a.db.QueryRow(readStmt, id)
	var nullDesc sql.NullString
	err := row.Scan(&m.Id, &m.Name, &nullDesc, &m.Saldo, &groupId, &m.NfcChipId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	m.Description = decodeNullableString(nullDesc)

	group, err := a.groups.Read(groupId)
	if err != nil {
		return nil, err
	}
	m.Group = group

	return m, nil
}

// Update saves the (changed) model in the database will return models.ErrGroupNotFound if group id is not associated with a group
func (a *AccountModel) Update(m *api.Account) error {
	updateStmt := `UPDATE accounts SET name=?, description=?, saldo=?, group_id=?, nfc_chip_uid=? WHERE id=?`

	_, err := a.db.Exec(updateStmt, m.Name, m.Description, m.Saldo, m.Group.Id, m.NfcChipId, m.Id)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1452 {
				return models.ErrGroupNotFound
			}
		}
		return err
	}

	return nil
}

// Delete deletes a account
func (a *AccountModel) Delete(id int32) error {

	deleteStmt := `DELETE FROM accounts WHERE id=?`

	_, err := a.db.Exec(deleteStmt, id)

	if err == sql.ErrNoRows {
		return models.ErrNotFound
	}

	return err
}

// UpdateSaldo provides a simpler update method for the saldo field
func (a *AccountModel) UpdateSaldo(id int32, newSaldo float64) error {
	_, err := a.db.Exec("UPDATE accounts SET saldo=? WHERE id=?", newSaldo, id)

	return err
}

// GetAll returns slice with all accounts in the database
func (a *AccountModel) GetAll() (*api.Accounts, error) {
	rows, err := a.db.Query("SELECT id, name, description, saldo, group_id FROM accounts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts, err := a.scanRowsToAccounts(rows)

	if err != nil {
		return nil, err
	}

	return &api.Accounts{Accounts: accounts}, nil
}

// GetAllByGroup returns slice with all accounts with given group id
func (a *AccountModel) GetAllByGroup(groupId int) ([]*api.Account, error) {
	rows, err := a.db.Query("SELECT id, name, description, saldo, group_id FROM accounts WHERE group_id=?", groupId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return a.scanRowsToAccounts(rows)
}

// scanRowsToAccounts returns slice of Accounts from given sql.Rows
func (a *AccountModel) scanRowsToAccounts(rows *sql.Rows) ([]*api.Account, error) {
	var accounts []*api.Account
	var groupIds []int32

	for rows.Next() {
		s := &api.Account{Group: &api.Group{}}

		var nullDesc sql.NullString

		err := rows.Scan(&s.Id, &s.Name, &nullDesc, &s.Saldo, &s.Group.Id)
		if err != nil {
			return nil, err
		}

		s.Description = decodeNullableString(nullDesc)

		groupIds = append(groupIds, s.Group.Id)
		accounts = append(accounts, s)
	}

	groups, err := a.groups.GetAllByIds(groupIds)

	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		account.Group = groups[account.Group.Id]
	}
	return accounts, nil
}
