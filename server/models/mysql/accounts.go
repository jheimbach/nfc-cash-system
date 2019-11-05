package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

type AccountModel struct {
	db *sql.DB
}

func (a *AccountModel) Create(name, description string, startSaldo float64, groupId int) error {
	createStmt := `INSERT INTO accounts (name, description, saldo, group_id) VALUES (?,?,?,?)`

	_, err := a.db.Exec(createStmt, name, description, startSaldo, groupId)

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

func (a *AccountModel) Read(id int) (*models.Account, error) {
	readStmt := `SELECT a.id, a.name, a.description, saldo, g.id, g.name, g.description
				 FROM accounts a
				     LEFT JOIN account_groups g on a.group_id = g.id
				 WHERE a.id = ?`

	m := &models.Account{
		Group: models.Group{},
	}
	row := a.db.QueryRow(readStmt, id)
	err := row.Scan(&m.ID, &m.Name, &m.Description, &m.Saldo, &m.Group.ID, &m.Group.Name, &m.Group.Description)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *AccountModel) Update(m *models.Account) error {
	updateStmt := `UPDATE accounts SET name=?, description=?, saldo=?, group_id=? WHERE id=?`

	_, err := a.db.Exec(updateStmt, m.Name, m.Description, m.Saldo, m.Group.ID, m.ID)

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

	return scanAccountFromRows(rows)
}

func (a *AccountModel) GetAllByGroup(groupId int) ([]*models.Account, error) {
	rows, err := a.db.Query("SELECT id, name, description, saldo, group_id FROM accounts WHERE group_id=?", groupId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return scanAccountFromRows(rows)
}

func (a *AccountModel) GetAllWithPaging(page, size int) (*models.AccountPaging, error) {
	getStmt := `SELECT id, name,description,saldo,group_id FROM accounts ORDER BY id DESC LIMIT ? OFFSET ?`
	offset := (page - 1) * size
	rows, err := a.db.Query(getStmt, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts, err := scanAccountFromRows(rows)
	if err != nil {
		return nil, err
	}

	var count int
	err = a.db.QueryRow("SELECT COUNT(id) FROM accounts").Scan(&count)
	if err != nil {
		return nil, err
	}

	return &models.AccountPaging{
		CurrentPage: page,
		MaxPage:     (count + size - 1) / size,
		Accounts:    accounts,
	}, nil
}

func scanAccountFromRows(rows *sql.Rows) ([]*models.Account, error) {
	var accounts []*models.Account

	for rows.Next() {
		s := &models.Account{Group: models.Group{}}

		err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Saldo, &s.Group.ID)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, s)
	}

	return accounts, nil
}
