package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

type GroupModel struct {
	db *sql.DB
}

func (g *GroupModel) Create(name, description string, canOverdraw bool) error {
	var nullDescription sql.NullString

	if description != "" {
		err := nullDescription.Scan(description)
		if err != nil {
			return err
		}
	}

	createStmt := "INSERT INTO `account_groups` (name, description, can_overdraw) VALUES (?,?,?)"
	_, err := g.db.Exec(createStmt, name, nullDescription, canOverdraw)
	return err
}

func (g *GroupModel) Read(id int) (*models.Group, error) {
	readStmt := "SELECT id, name, description, can_overdraw FROM `account_groups` WHERE id = ?"

	var group models.Group
	row := g.db.QueryRow(readStmt, id)

	var nullDesc sql.NullString
	err := row.Scan(&group.ID, &group.Name, &nullDesc, &group.CanOverDraw)

	if nullDesc.Valid {
		group.Description = nullDesc.String
	}

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

	res, err := g.db.Exec(
		"UPDATE `account_groups` SET name=?,description=?, can_overdraw=? WHERE id=?",
		group.Name,
		group.Description,
		group.CanOverDraw,
		group.ID,
	)

	if err != nil {
		return nil, err
	}

	// check whether update has affected some rows, if not return models.ErrNotFound
	if affected, err := res.RowsAffected(); affected == 0 || err != nil {
		return nil, models.ErrNotFound
	}

	return &group, nil
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
