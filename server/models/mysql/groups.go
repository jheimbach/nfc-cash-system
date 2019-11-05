package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

type GroupModel struct {
	db *sql.DB
}

func (g *GroupModel) Create(name, description string) error {
	createStmt := "INSERT INTO `account_groups` (name, description) VALUES (?,?)"
	_, err := g.db.Exec(createStmt, name, description)
	return err
}

func (g *GroupModel) Read(id int) (*models.Group, error) {
	readStmt := "SELECT * FROM `account_groups` WHERE id = ?"

	var group models.Group
	row := g.db.QueryRow(readStmt, id)

	err := row.Scan(&group.ID, &group.Name, &group.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return &group, nil
}

func (g *GroupModel) Update(group models.Group) (*models.Group, error) {
	if group.ID == 0 {
		return nil, models.ErrModelNotSaved
	}

	current, err := g.Read(group.ID)
	if err != nil {
		return nil, err
	}

	if current.Name != group.Name {
		current.Name = group.Name
	}

	if current.Description != group.Description {
		current.Description = group.Description
	}

	_, err = g.db.Exec(
		"UPDATE `account_groups` SET name=?,description=? WHERE id=?",
		current.Name,
		current.Description,
		current.ID,
	)

	if err != nil {
		return nil, err
	}

	return current, nil
}

func (g *GroupModel) Delete(id int) error {
	_, err := g.db.Exec("DELETE FROM `account_groups` WHERE id=?", id)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1451 {
				return models.ErrNonEmptyDelete
			}
		}

		return err
	}

	return nil
}
