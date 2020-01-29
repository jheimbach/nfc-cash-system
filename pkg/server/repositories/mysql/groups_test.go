package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jheimbach/nfc-cash-system/api"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test"
	"github.com/jheimbach/nfc-cash-system/pkg/server/repositories"
	isPkg "github.com/matryer/is"
)

func TestGroupModel_Create(t *testing.T) {
	is := isPkg.New(t)

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
			defer teardownDB(_conn)()

			got, err := _groupModel.Create(context.Background(), tt.args.name, tt.args.description, tt.args.canOverdraw)
			is.NoErr(err)
			is.Equal(got, &tt.want) // does not return expected group

			var dbGroup api.Group
			var nullDesc sql.NullString
			row := _conn.QueryRow("SELECT id, name, description,can_overdraw FROM `account_groups` WHERE id = ?", tt.want.Id)
			err = row.Scan(&dbGroup.Id, &dbGroup.Name, &nullDesc, &dbGroup.CanOverdraw)
			is.NoErr(err)

			dbGroup.Description = decodeNullableString(nullDesc)

			is.Equal(dbGroup, tt.want)
		})
	}
}

func TestGroupModel_Read(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)

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
			expectedErr: repositories.ErrNotFound,
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
			defer teardownDB(_conn)()

			is := is.New(t)
			if tt.insertGroup {
				err := insertMockGroup(t, tt.want)
				is.NoErr(err) //could not insert test group
			}

			got, err := _groupModel.Read(context.Background(), tt.want.Id)

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
	test.IsIntegrationTest(t)

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
			expectedErr: repositories.ErrModelNotSaved,
		},
		{
			name: "group without id will not be updated",
			want: &api.Group{
				Name:        "testgroup",
				Description: "test description",
			},
			wantErr:     true,
			expectedErr: repositories.ErrModelNotSaved,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownDB(_conn)()

			is := isPkg.New(t)
			if !tt.skipInsert {
				err := insertMockGroup(t, tt.want)
				is.NoErr(err) //could not insert test group
			}

			got, err := _groupModel.Update(context.Background(), tt.want)

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
	test.IsIntegrationTest(t)

	tests := []struct {
		name     string
		group    api.Group
		accounts []api.Account
		wantErr  error
	}{
		{
			name: "empty group delete",
			group: api.Group{
				Id:   1,
				Name: "test",
			},
		},
		{
			name: "trying to delete nonempty group, return err",
			group: api.Group{
				Id:   1,
				Name: "test",
			},
			accounts: []api.Account{
				{
					Name:      "test",
					NfcChipId: "testchipid",
					Group:     &api.Group{Id: 1},
				},
			},
			wantErr: repositories.ErrNonEmptyDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownDB(_conn)()
			is := isPkg.New(t)

			err := insertMockGroup(t, &tt.group)
			is.NoErr(err) // could not create test group

			if tt.accounts != nil {
				for _, account := range tt.accounts {
					err := insertTestAccount(t, account)
					is.NoErr(err) // could not create test account
				}
			}

			err = _groupModel.Delete(context.Background(), tt.group.Id)
			if tt.wantErr != nil {
				if tt.wantErr != err {
					t.Errorf("got err %q, expected %q", err, tt.wantErr)
				}
				return
			}

			is.NoErr(err) // could not delete group

			// check if group exists
			var groupCount int
			err = _conn.QueryRow("SELECT COUNT(*) from `account_groups` WHERE id=?", int(tt.group.Id)).Scan(&groupCount)
			is.Equal(groupCount, 0) // group still exists
		})
	}
}

func TestGroupModel_GetAll(t *testing.T) {
	test.IsIntegrationTest(t)

	is := isPkg.New(t)

	teardown := initDBForGroupList(t)
	defer teardown()

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
			want:      wantGroupList(0, 0),
			wantCount: 10,
		},
		{
			name: "get groups with limit",
			input: args{
				limit:  5,
				offset: 0,
			},
			want:      wantGroupList(5, 0),
			wantCount: 10,
		},
		{
			name: "get groups with limit and offset",
			input: args{
				limit:  5,
				offset: 5,
			},
			want:      wantGroupList(5, 5),
			wantCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got, count, err := _groupModel.GetAll(context.Background(), tt.input.limit, tt.input.offset)
			is.NoErr(err)

			is.Equal(got, tt.want)
			is.Equal(count, tt.wantCount)
		})
	}
}

func TestGroupModel_GetAllByIds(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)

	teardown := initDBForGroupList(t)
	defer teardown()

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

			got, err := _groupModel.GetAllByIds(context.Background(), tt.input)
			if tt.wantErr != nil {
				is.Equal(err, tt.wantErr)
				return
			}

			is.NoErr(err)

			is.Equal(got, tt.want)
		})
	}
}

func initDBForGroupList(t *testing.T) func() error {
	t.Helper()

	err := test.SetupDB(_conn, dataFor("group_list"))
	if err != nil {
		t.Fatal(err)
	}

	return teardownDB(_conn)
}

func insertMockGroup(t *testing.T, group *api.Group) error {
	t.Helper()

	_, err := _conn.Exec(
		"INSERT INTO `account_groups` (id, name, description, can_overdraw) VALUES (?,?,?,?)",
		group.Id, group.Name,
		createNullableString(group.Description),
		group.CanOverdraw,
	)
	return err
}

func wantGroupList(limit, offset int32) []*api.Group {
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
