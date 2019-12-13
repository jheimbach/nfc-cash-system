package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
)

func TestGroupModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		name, description string
		canOverdraw       bool
	}
	tests := []struct {
		name string
		args args
		want api.Group
	}{
		{
			name: "create group without description",
			args: args{
				name: "testgroup",
			},
			want: api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "",
				CanOverdraw: false,
			},
		},
		{
			name: "create group with description",
			args: args{
				name:        "testgroup",
				description: "with description",
			},
			want: api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "with description",
			},
		},
		{
			name: "create group with canoverdraw",
			args: args{
				name:        "testgroup",
				canOverdraw: true,
			},
			want: api.Group{
				Id:          1,
				Name:        "testgroup",
				CanOverdraw: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dbSetup()
			defer dbTeardown()

			model := GroupModel{
				db: db,
			}

			got, err := model.Create(context.Background(), tt.args.name, tt.args.description, tt.args.canOverdraw)
			is.NoErr(err)
			is.Equal(got, &tt.want) // does not return expected group

			var dbGroup api.Group
			var nullDesc sql.NullString
			row := db.QueryRow("SELECT id, name, description,can_overdraw FROM `account_groups` WHERE id = ?", tt.want.Id)
			err = row.Scan(&dbGroup.Id, &dbGroup.Name, &nullDesc, &dbGroup.CanOverdraw)
			is.NoErr(err)

			dbGroup.Description = decodeNullableString(nullDesc)

			is.Equal(dbGroup, tt.want)
		})
	}
}

func TestGroupModel_Read(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	tests := []struct {
		name        string
		want        *api.Group
		insertGroup bool
		wantErr     bool
		expectedErr error
	}{
		{
			name: "load group without description",
			want: &api.Group{
				Id:   1,
				Name: "testgroup",
			},
			insertGroup: true,
		},
		{
			name: "load group with description",
			want: &api.Group{
				Id:          2,
				Name:        "testgroup",
				Description: "with description",
			},
			insertGroup: true,
		},
		{
			name: "if group id does not exist, return models.ErrNotFound",
			want: &api.Group{
				Id: 3,
			},
			insertGroup: false,
			wantErr:     true,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "load group with Canoverdraw",
			want: &api.Group{
				Id:          4,
				Name:        "testgroup4",
				CanOverdraw: true,
			},
			insertGroup: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSetup()
			defer dbTeardown()

			model := GroupModel{
				db: db,
			}

			is := is.New(t)
			if tt.insertGroup {
				_, err := db.Exec("INSERT INTO `account_groups` (id, name, description, can_overdraw) VALUES (?,?,?,?)", tt.want.Id, tt.want.Name, createNullableString(tt.want.Description), tt.want.CanOverdraw)
				is.NoErr(err)
			}

			got, err := model.Read(context.Background(), tt.want.Id)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}
			is.NoErr(err)

			is.Equal(got, tt.want)
		})
	}
}

func TestGroupModel_Update(t *testing.T) {
	isIntegrationTest(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	tests := []struct {
		name        string
		insert      api.Group
		want        *api.Group
		wantErr     bool
		expectedErr error
		skipInsert  bool
	}{
		{
			name: "update group description",
			insert: api.Group{
				Name: "testgroup",
			},
			want: &api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "test description",
			},
		},
		{
			name: "update group name",
			insert: api.Group{
				Name:        "testgroup",
				Description: "non empty",
			},
			want: &api.Group{
				Id:          1,
				Name:        "test",
				Description: "non empty",
			},
		},
		{
			name: "update nothing",
			insert: api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "non empty",
			},
			want: &api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "non empty",
			},
		},
		{
			name: "update can_overdraw",
			insert: api.Group{
				Name:        "testgroup",
				CanOverdraw: false,
			},
			want: &api.Group{
				Id:          1,
				Name:        "testgroup",
				CanOverdraw: true,
			},
		},
		{
			name:        "empty group will not be updated",
			want:        &api.Group{},
			skipInsert:  true,
			wantErr:     true,
			expectedErr: models.ErrModelNotSaved,
		},
		{
			name: "group without id will not be updated",
			want: &api.Group{
				Name:        "testgroup",
				Description: "test description",
			},
			wantErr:     true,
			expectedErr: models.ErrModelNotSaved,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSetup()
			defer dbTeardown()

			model := GroupModel{
				db: db,
			}

			is := isPkg.New(t)
			if !tt.skipInsert {
				_, err := db.Exec("INSERT INTO `account_groups` (name, description,can_overdraw) VALUES (?,?,?)", tt.insert.Name, createNullableString(tt.insert.Description), tt.insert.CanOverdraw)
				is.NoErr(err)
			}

			got, err := model.Update(context.Background(), tt.want)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}

			is.NoErr(err)

			is.Equal(got, tt.want)

		})
	}
}

func TestGroupModel_Delete(t *testing.T) {
	isIntegrationTest(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("empty group delete", func(t *testing.T) {
		is := isPkg.New(t)
		dbSetup()
		defer dbTeardown()

		model := GroupModel{
			db: db,
		}
		res, err := db.Exec("INSERT INTO `account_groups` (name) VALUES (?)", "test")
		is.NoErr(err)

		groupId, err := res.LastInsertId()
		is.NoErr(err)

		err = model.Delete(context.Background(), int32(groupId))
		is.NoErr(err)

		var groupName string
		err = db.QueryRow("SELECT name from `account_groups` WHERE id=?", int(groupId)).Scan(&groupName)
		if err == nil {
			t.Errorf("wanted err, got none")
		}
		if err != sql.ErrNoRows {
			t.Errorf("got %v but wanted %v err", sql.ErrNoRows, err)
		}
	})

	t.Run("trying to delete nonempty group, return err", func(t *testing.T) {
		is := isPkg.New(t)
		dbSetup()
		defer dbTeardown()

		model := GroupModel{
			db: db,
		}
		res, err := db.Exec("INSERT INTO `account_groups` (name) VALUES (?)", "test")
		is.NoErr(err)

		groupId, err := res.LastInsertId()
		is.NoErr(err)

		_, err = db.Exec("INSERT INTO `accounts` (name, group_id, nfc_chip_uid) VALUES (?,?,?)", "test", int(groupId), "testchipid")
		is.NoErr(err)

		err = model.Delete(context.Background(), int32(groupId))
		if err == nil {
			t.Errorf("expected error, got none")
		}
		if err != models.ErrNonEmptyDelete {
			t.Errorf("got %v, wanted %v", err, models.ErrNonEmptyDelete)
		}
	})
}

func TestGroupModel_GetAll(t *testing.T) {
	isIntegrationTest(t)

	is := isPkg.New(t)
	db, dbTeardown := dbInitializedForGroupList(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := GroupModel{
		db: db,
	}

	type args struct {
		limit, offset int32
	}
	tests := []struct {
		name      string
		input     args
		want      []*api.Group
		wantCount int
		wantErr   error
	}{
		{
			name: "get all groups",
			input: args{
				limit:  0,
				offset: 0,
			},
			want:      groupList(0, 0),
			wantCount: 10,
		},
		{
			name: "get groups with limit",
			input: args{
				limit:  5,
				offset: 0,
			},
			want:      groupList(5, 0),
			wantCount: 10,
		},
		{
			name: "get groups with limit and offset",
			input: args{
				limit:  5,
				offset: 5,
			},
			want:      groupList(5, 5),
			wantCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got, count, err := model.GetAll(context.Background(), tt.input.limit, tt.input.offset)
			is.NoErr(err)

			is.Equal(got, tt.want)
			is.Equal(count, tt.wantCount)
		})
	}
}

func TestGroupModel_GetAllByIds(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbTeardown := dbInitializedForGroupList(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	tests := []struct {
		name    string
		input   []int32
		want    map[int32]*api.Group
		wantErr error
	}{
		{
			name:  "get single group with id 1",
			input: []int32{1},
			want:  map[int32]*api.Group{1: mockGroupOne},
		},
		{
			name:  "get multiple groups with id 1 and 2",
			input: []int32{1, 2},
			want:  mockGroupMap,
		},
		{
			name:  "input ids are not found",
			input: []int32{100, 200},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			model := GroupModel{
				db: db,
			}

			got, err := model.GetAllByIds(context.Background(), tt.input)
			if tt.wantErr != nil {
				is.Equal(err, tt.wantErr)
				return
			}

			is.NoErr(err)

			is.Equal(got, tt.want)
		})
	}
}

func dbInitializedForGroupList(t *testing.T) (*sql.DB, func()) {
	db, setup, teardown := getTestDb(t)
	setup("../testdata/group_list.sql")

	return db, teardown
}

func groupList(limit, offset int32) []*api.Group {
	groups := []*api.Group{
		{
			Id:   1,
			Name: "testgroup1",
		},
		{
			Id:          2,
			Name:        "testgroup2",
			Description: "with description",
		},
		{
			Id:          3,
			Name:        "testgroup3",
			CanOverdraw: true,
		},
		{
			Id:          4,
			Name:        "testgroup4",
			Description: "with description",
			CanOverdraw: true,
		},
		{
			Id:          5,
			Name:        "testgroup5",
			CanOverdraw: true,
		},
		{
			Id:          6,
			Name:        "testgroup6",
			Description: "with description",
			CanOverdraw: true,
		},
		{
			Id:   7,
			Name: "testgroup7",
		},
		{
			Id:          8,
			Name:        "testgroup8",
			Description: "with description",
		},
		{
			Id:   9,
			Name: "testgroup9",
		},
		{
			Id:          10,
			Name:        "testgroup10",
			Description: "with description",
		},
	}

	if limit > 0 {
		var off int32
		if offset > 0 {
			off = offset
		}
		return groups[off : off+limit]
	}

	return groups
}
