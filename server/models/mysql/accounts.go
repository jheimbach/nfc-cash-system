package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

type AccountModel struct {
	db *sql.DB
}

func (a *AccountModel) Create(name, description string, startSaldo float64, groupId int, nfcChipId string) error {
	nullDescription := createNullableString(description)

	createStmt := `INSERT INTO accounts (name, description, saldo, group_id, nfc_chip_uid) VALUES (?,?,?,?,?)`

	_, err := a.db.Exec(createStmt, name, nullDescription, startSaldo, groupId, nfcChipId)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1452 {
				return models.ErrGroupNotFound
			}
			if err.Number == 1062 {
				return models.ErrDuplicateNfcChipId
			}
		}
		return err
	}

	return nil
}

func (a *AccountModel) Read(id int) (*models.Account, error) {
	readStmt := `SELECT a.id, a.name, a.description, a.saldo, a.nfc_chip_uid, g.id, g.name, g.description
				 FROM accounts a
				     LEFT JOIN account_groups g on a.group_id = g.id
				 WHERE a.id = ?`

	m := &models.Account{
		Group: &models.Group{},
	}
	row := a.db.QueryRow(readStmt, id)
	var nullDesc sql.NullString
	var groupNullDesc sql.NullString
	err := row.Scan(&m.ID, &m.Name, &nullDesc, &m.Saldo, &m.NfcChipId, &m.Group.ID, &m.Group.Name, &groupNullDesc)
	if err != nil {
		return nil, err
	}
	m.Description = decodeNullableString(nullDesc)
	m.Group.Description = decodeNullableString(groupNullDesc)

	return m, nil
}

func (a *AccountModel) Update(m *models.Account) error {
	updateStmt := `UPDATE accounts SET name=?, description=?, saldo=?, group_id=?, nfc_chip_uid=? WHERE id=?`

	_, err := a.db.Exec(updateStmt, m.Name, m.Description, m.Saldo, m.Group.ID, m.NfcChipId, m.ID)

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

func (a *AccountModel) Delete(id int) error {

	deleteStmt := `DELETE FROM accounts WHERE id=?`

	_, err := a.db.Exec(deleteStmt, id)

	if err == sql.ErrNoRows {
		return models.ErrNotFound
	}

	return err
}

func (a *AccountModel) UpdateSaldo(id int, newSaldo float64) error {
	_, err := a.db.Exec("UPDATE accounts SET saldo=? WHERE id=?", newSaldo, id)

	return err
}

func (a *AccountModel) GetAll() ([]*models.Account, error) {
	rows, err := a.db.Query("SELECT id, name, description, saldo, group_id FROM accounts")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return scanRowsToAccounts(rows)
}

func (a *AccountModel) GetAllByGroup(groupId int) ([]*models.Account, error) {
	rows, err := a.db.Query("SELECT id, name, description, saldo, group_id FROM accounts WHERE group_id=?", groupId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return scanRowsToAccounts(rows)
}

func (a *AccountModel) GetAllWithPaging(page, size int) (*models.AccountPaging, error) {
	getStmt := `SELECT id, name,description,saldo,group_id FROM accounts ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := a.db.Query(getStmt, size, pageOffset(page, size))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts, err := scanRowsToAccounts(rows)
	if err != nil {
		return nil, err
	}

	count, err := countAllIds(a.db, "SELECT COUNT(id) FROM accounts")
	if err != nil {
		return nil, err
	}

	return &models.AccountPaging{
		CurrentPage: page,
		MaxPage:     maxPageCount(count, size),
		Accounts:    accounts,
	}, nil
}

func scanRowsToAccounts(rows *sql.Rows) ([]*models.Account, error) {
	var accounts []*models.Account

	for rows.Next() {
		s := &models.Account{Group: &models.Group{}}

		var nullDesc sql.NullString

		err := rows.Scan(&s.ID, &s.Name, &nullDesc, &s.Saldo, &s.Group.ID)
		if err != nil {
			return nil, err
		}

		s.Description = decodeNullableString(nullDesc)

		accounts = append(accounts, s)
	}

	return accounts, nil
}
