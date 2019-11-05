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
