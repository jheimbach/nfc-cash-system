package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
	"testing"
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
		want models.Group
	}{
		{
			name: "create group without description",
			args: args{
				name: "testgroup",
			},
			want: models.Group{
				ID:          1,
				Name:        "testgroup",
				Description: "",
				CanOverDraw: false,
			},
		},
		{
			name: "create group with description",
			args: args{
				name:        "testgroup",
				description: "with description",
			},
			want: models.Group{
				ID:          1,
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
			want: models.Group{
				ID:          1,
				Name:        "testgroup",
				CanOverDraw: true,
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

			err := model.Create(tt.args.name, tt.args.description, tt.args.canOverdraw)
			is.NoErr(err)

			var got models.Group
			var nullDesc sql.NullString
			row := db.QueryRow("SELECT id, name, description,can_overdraw FROM `account_groups` WHERE id = ?", tt.want.ID)
			err = row.Scan(&got.ID, &got.Name, &nullDesc, &got.CanOverDraw)
			is.NoErr(err)

			got.Description = decodeNullableString(nullDesc)

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
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
		want        *models.Group
		insertGroup bool
		wantErr     bool
		expectedErr error
	}{
		{
			name: "load group without description",
			want: &models.Group{
				ID:   1,
				Name: "testgroup",
			},
			insertGroup: true,
		},
		{
			name: "load group with description",
			want: &models.Group{
				ID:          2,
				Name:        "testgroup",
				Description: "with description",
			},
			insertGroup: true,
		},
		{
			name: "if group id does not exist, return models.ErrNotFound",
			want: &models.Group{
				ID: 3,
			},
			insertGroup: false,
			wantErr:     true,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "load group with Canoverdraw",
			want: &models.Group{
				ID:          4,
				Name:        "testgroup4",
				CanOverDraw: true,
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
				_, err := db.Exec("INSERT INTO `account_groups` (id, name, description, can_overdraw) VALUES (?,?,?,?)", tt.want.ID, tt.want.Name, createNullableString(tt.want.Description), tt.want.CanOverDraw)
				is.NoErr(err)
			}

			got, err := model.Read(tt.want.ID)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}
			is.NoErr(err)

			if *got != *tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestGroupModel_Update(t *testing.T) {
	isIntegrationTest(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	tests := []struct {
		name        string
		insert      models.Group
		want        models.Group
		wantErr     bool
		expectedErr error
		skipInsert  bool
	}{
		{
			name: "update group description",
			insert: models.Group{
				Name: "testgroup",
			},
			want: models.Group{
				ID:          1,
				Name:        "testgroup",
				Description: "test description",
			},
		},
		{
			name: "update group name",
			insert: models.Group{
				Name:        "testgroup",
				Description: "non empty",
			},
			want: models.Group{
				ID:          1,
				Name:        "test",
				Description: "non empty",
			},
		},
		{
			name: "update can_overdraw",
			insert: models.Group{
				Name:        "testgroup",
				CanOverDraw: false,
			},
			want: models.Group{
				ID:          1,
				Name:        "testgroup",
				CanOverDraw: true,
			},
		},
		{
			name:        "empty group will not be updated",
			want:        models.Group{},
			skipInsert:  true,
			wantErr:     true,
			expectedErr: models.ErrModelNotSaved,
		},
		{
			name: "group with non existent id returns models.ErrNotFound",
			want: models.Group{
				ID: 12,
			},
			skipInsert:  true,
			wantErr:     true,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "group without id will not be updated",
			want: models.Group{
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
				_, err := db.Exec("INSERT INTO `account_groups` (name, description,can_overdraw) VALUES (?,?,?)", tt.insert.Name, createNullableString(tt.insert.Description), tt.insert.CanOverDraw)
				is.NoErr(err)
			}

			got, err := model.Update(tt.want)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}

			is.NoErr(err)

			if *got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
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

		err = model.Delete(int(groupId))
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

		err = model.Delete(int(groupId))
		if err == nil {
			t.Errorf("expected error, got none")
		}
		if err != models.ErrNonEmptyDelete {
			t.Errorf("got %v, wanted %v", err, models.ErrNonEmptyDelete)
		}
	})
}
