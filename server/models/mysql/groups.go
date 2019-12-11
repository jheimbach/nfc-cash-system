package mysql

import (
	"database/sql"
	"strings"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/go-sql-driver/mysql"
)

// GroupModel provides API for the account_groups table
type GroupModel struct {
	db *sql.DB
}

func NewGroupModel(db *sql.DB) *GroupModel {
	return &GroupModel{db: db}
}

// Creates inserts new group with given fields
func (g *GroupModel) Create(name, description string, canOverdraw bool) (*api.Group, error) {
	nullDescription := createNullableString(description)

	createStmt := "INSERT INTO `account_groups` (name, description, can_overdraw) VALUES (?,?,?)"
	res, err := g.db.Exec(createStmt, name, nullDescription, canOverdraw)

	if err != nil {
		return nil, err
	}

	group := &api.Group{
		Name:        name,
		Description: description,
		CanOverdraw: canOverdraw,
	}

	// mysql returns always nil as error value on LastInsertId(), we don't have to check it
	lastId, _ := res.LastInsertId()

	group.Id = int32(lastId)

	return group, nil
}

// Read returns models.Group struct for given id, will return models.ErrNotFound if no group is found
func (g *GroupModel) Read(id int32) (*api.Group, error) {
	readStmt := "SELECT id, name, description, can_overdraw FROM `account_groups` WHERE id = ?"

	var group api.Group
	row := g.db.QueryRow(readStmt, id)

	var nullDesc sql.NullString
	err := row.Scan(&group.Id, &group.Name, &nullDesc, &group.CanOverdraw)
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
func (g *GroupModel) Update(group *api.Group) (*api.Group, error) {
	if group.Id == 0 {
		return nil, models.ErrModelNotSaved
	}

	_, err := g.db.Exec(
		"UPDATE `account_groups` SET name=?,description=?, can_overdraw=? WHERE id=?",
		group.Name,
		group.Description,
		group.CanOverdraw,
		group.Id,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}

// Delete removes group with given id from the database
// returns models.ErrNonEmptyDelete if accounts are associated with group
func (g *GroupModel) Delete(id int32) error {
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

func (g *GroupModel) GetAll() (*api.Groups, error) {

	rows, err := g.db.Query("SELECT id, name, description, can_overdraw FROM account_groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups, err := scanGroups(rows)
	if err != nil {
		return nil, err
	}

	return &api.Groups{Groups: groups}, nil
}

func (g *GroupModel) GetAllByIds(ids []int32) (map[int32]*api.Group, error) {
	if len(ids) == 0 {
		return nil, models.ErrNotFound //todo maybe own error?
	}

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	readStmt := `SELECT id,name,description,can_overdraw FROM account_groups WHERE id IN (?` + strings.Repeat(",?", len(ids)-1) + `)`
	rows, err := g.db.Query(readStmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups, err := scanGroups(rows)
	if err != nil {
		return nil, err
	}

	if groups == nil {
		return nil, nil
	}

	m := make(map[int32]*api.Group, len(groups))

	for _, group := range groups {
		m[group.Id] = group
	}

	return m, nil
}

func scanGroups(rows *sql.Rows) ([]*api.Group, error) {
	var groups []*api.Group

	for rows.Next() {
		g := &api.Group{}

		var descriptionNullable sql.NullString
		err := rows.Scan(&g.Id, &g.Name, &descriptionNullable, &g.CanOverdraw)
		g.Description = decodeNullableString(descriptionNullable)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
