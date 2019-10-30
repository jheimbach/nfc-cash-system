package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/pkg/models"
	"testing"
)

func TestGroupModel_Create(t *testing.T) {
	db, teardown := getTestDb(t)
	defer teardown()

	model := GroupModel{
		db: db,
	}

	t.Run("create group without description", func(t *testing.T) {
		err := model.Create("testgroup", "")
		assertErrIsNil(t, err, false)

		want := models.Group{
			ID:          1,
			Name:        "testgroup",
			Description: "",
		}

		var got models.Group

		row := db.QueryRow("SELECT id, name, description FROM `groups` WHERE id = ?", want.ID)
		err = row.Scan(&got.ID, &got.Name, &got.Description)
		assertErrIsNil(t, err, false)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("create group with description", func(t *testing.T) {

		want := models.Group{
			ID:          2,
			Name:        "testgroup",
			Description: "with description",
		}
		err := model.Create(want.Name, want.Description)
		assertErrIsNil(t, err, false)

		var got models.Group

		row := db.QueryRow("SELECT id, name, description FROM `groups` WHERE id = ?", want.ID)
		err = row.Scan(&got.ID, &got.Name, &got.Description)
		assertErrIsNil(t, err, false)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestGroupModel_Read(t *testing.T) {
	db, teardown := getTestDb(t)
	defer teardown()

	model := GroupModel{
		db: db,
	}

	tests := []struct {
		name        string
		want        *models.Group
		insertGroup bool
		wantErr     bool
		expectedErr error
	}{
		{
			name: "insert group with no description",
			want: &models.Group{
				ID:          1,
				Name:        "testgroup",
				Description: "",
			},
			insertGroup: true,
		},
		{
			name: "insert group with description",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.insertGroup {
				_, err := db.Exec("INSERT INTO `groups` (name, description) VALUES (?,?)", tt.want.Name, tt.want.Description)
				assertErrIsNil(t, err, false)
			}

			got, err := model.Read(tt.want.ID)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}
			assertErrIsNil(t, err, true)

			if *got != *tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestGroupModel_Update(t *testing.T) {
	db, teardown := getTestDb(t)
	defer teardown()

	model := GroupModel{
		db: db,
	}

	tests := []struct {
		name        string
		insert      models.Group
		want        models.Group
		wantErr     bool
		expectedErr error
	}{
		{
			name: "insert group with no description, add description and return group",
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
			name: "insert group and change name",
			insert: models.Group{
				Name:        "testgroup",
				Description: "non empty",
			},
			want: models.Group{
				ID:          2,
				Name:        "test",
				Description: "non empty",
			},
		},
		{
			name:        "empty group will not be updated",
			want:        models.Group{},
			wantErr:     true,
			expectedErr: models.ErrModelNotSaved,
		},
		{
			name: "group with non existend id returns models.ErrNotFound",
			want: models.Group{
				ID: 12,
			},
			wantErr:     true,
			expectedErr: models.ErrNotFound,
		},
		{
			name: "group without a id will not be updated",
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
			if tt.insert.Name != "" {
				_, err := db.Exec("INSERT INTO `groups` (name, description) VALUES (?,?)", tt.insert.Name, tt.insert.Description)
				assertErrIsNil(t, err, true)
			}

			got, err := model.Update(tt.want)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v but expected %v", err, tt.expectedErr)
				}
				return
			}

			assertErrIsNil(t, err, true)

			if *got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestGroupModel_Delete(t *testing.T) {
	t.Run("empty group delete", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := GroupModel{
			db: db,
		}
		res, err := db.Exec("INSERT INTO `groups` (name) VALUES (?)", "test")
		assertErrIsNil(t, err, true)

		groupId, err := res.LastInsertId()
		assertErrIsNil(t, err, true)

		err = model.Delete(int(groupId))
		assertErrIsNil(t, err, false)

		var groupName string
		err = db.QueryRow("SELECT name from `groups` WHERE id=?", int(groupId)).Scan(&groupName)
		if err == nil {
			t.Errorf("wanted err, got none")
		}
		if err != sql.ErrNoRows {
			t.Errorf("got %v but wanted %v err", sql.ErrNoRows, err)
		}
	})
	t.Run("trying to delete nonempty group, return err", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := GroupModel{
			db: db,
		}
		res, err := db.Exec("INSERT INTO `groups` (name) VALUES (?)", "test")
		assertErrIsNil(t, err, true)

		groupId, err := res.LastInsertId()
		assertErrIsNil(t, err, true)

		_, err = db.Exec("INSERT INTO `accounts` (name, group_id) VALUES (?,?)", "test", int(groupId))
		assertErrIsNil(t, err, true)

		err = model.Delete(int(groupId))
		if err == nil {
			t.Errorf("expected error, got none")
		}
		if err != models.ErrNonEmptyDelete {
			t.Errorf("got %v, wanted %v", err, models.ErrNonEmptyDelete)
		}
	})
}
