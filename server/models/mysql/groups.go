package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

// GroupModel provides API for the account_groups table
type GroupModel struct {
	db *sql.DB
}

// Creates inserts new group with given fields
func (g *GroupModel) Create(name, description string, canOverdraw bool) error {
	nullDescription := createNullableString(description)

	createStmt := "INSERT INTO `account_groups` (name, description, can_overdraw) VALUES (?,?,?)"
	_, err := g.db.Exec(createStmt, name, nullDescription, canOverdraw)
	return err
}

// Read returns models.Group struct for given id, will return models.ErrNotFound if no group is found
func (g *GroupModel) Read(id int) (*models.Group, error) {
	readStmt := "SELECT id, name, description, can_overdraw FROM `account_groups` WHERE id = ?"

	var group models.Group
	row := g.db.QueryRow(readStmt, id)

	var nullDesc sql.NullString
	err := row.Scan(&group.ID, &group.Name, &nullDesc, &group.CanOverDraw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	group.Description = decodeNullableString(nullDesc)

	return &group, nil
}

// Update saves given (changed) group to the database
// will return models.ErrNotFound if group is not found
// NOTE: every field will be overwritten with given value
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

// Delete removes group with given id from the database
// returns models.ErrNonEmptyDelete if accounts are associated with group
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
