package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jheimbach/nfc-cash-system/api"
	"github.com/jheimbach/nfc-cash-system/pkg/server/repositories"
)

// GroupRepository provides API for the account_groups table
type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// Creates inserts new group with given fields
func (g *GroupRepository) Create(ctx context.Context, name, description string, canOverdraw bool) (*api.Group, error) {
	nullDescription := createNullableString(description)

	createStmt := "INSERT INTO `account_groups` (name, description, can_overdraw) VALUES (?,?,?)"
	res, err := g.db.ExecContext(ctx, createStmt, name, nullDescription, canOverdraw)

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
func (g *GroupRepository) Read(ctx context.Context, id int32) (*api.Group, error) {
	readStmt := "SELECT id, name, description, can_overdraw FROM `account_groups` WHERE id = ?"

	var group api.Group
	row := g.db.QueryRowContext(ctx, readStmt, id)

	var nullDesc sql.NullString
	err := row.Scan(&group.Id, &group.Name, &nullDesc, &group.CanOverdraw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	group.Description = decodeNullableString(nullDesc)

	return &group, nil
}

// Update saves given (changed) group to the database
// will return models.ErrNotFound if group is not found
// NOTE: every field will be overwritten with given value
func (g *GroupRepository) Update(ctx context.Context, group *api.Group) (*api.Group, error) {
	if group.Id == 0 {
		return nil, repositories.ErrModelNotSaved
	}

	_, err := g.db.ExecContext(ctx,
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
func (g *GroupRepository) Delete(ctx context.Context, id int32) error {
	_, err := g.db.ExecContext(ctx, "DELETE FROM `account_groups` WHERE id=?", id)

	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1451 {
				return repositories.ErrNonEmptyDelete
			}
		}

		return err
	}

	return nil
}

func (g *GroupRepository) GetAll(ctx context.Context, limit, offset int32) ([]*api.Group, int, error) {
	stmt := "SELECT id, name, description, can_overdraw FROM account_groups"

	var args []interface{}
	if limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT ?", stmt)
		args = append(args, limit)
		if offset > 0 {
			stmt = fmt.Sprintf("%s OFFSET ?", stmt)
			args = append(args, offset)
		}
	}

	rows, err := g.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	groups, err := scanGroups(rows)
	if err != nil {
		return nil, 0, err
	}

	totalCount := len(groups)
	if limit > 0 {
		totalCount, err = g.countAll(ctx)
		if err != nil {
			return nil, 0, err
		}
	}

	return groups, totalCount, nil
}

func (g *GroupRepository) countAll(ctx context.Context) (int, error) {
	stmt := `SELECT COUNT(id) FROM account_groups`

	var totalCount int
	err := g.db.QueryRowContext(ctx, stmt).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (g *GroupRepository) GetAllByIds(ctx context.Context, ids []int32) (map[int32]*api.Group, error) {
	if len(ids) == 0 {
		return nil, repositories.ErrNotFound //todo maybe own error?
	}

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	readStmt := `SELECT id,name,description,can_overdraw FROM account_groups WHERE id IN (?` + strings.Repeat(",?", len(ids)-1) + `)`
	rows, err := g.db.QueryContext(ctx, readStmt, args...)

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
				return nil, repositories.ErrNotFound
			}
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
