package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
	"io/ioutil"
	"testing"
)

func TestAccountModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	type fields struct {
		name, description string
		saldo             float64
		groupId           int
	}
	tests := []struct {
		name        string
		input       fields
		want        fields
		wantErr     bool
		expectedErr error
	}{
		{
			name: "creates account",
			input: fields{
				name:        "tim",
				description: "",
				saldo:       12,
				groupId:     1,
			},
			want: fields{
				name:        "tim",
				description: "",
				saldo:       12,
				groupId:     1,
			},
		},
		{
			name: "creates account but group does not exists",
			input: fields{
				name:        "tim",
				description: "",
				saldo:       12,
				groupId:     100,
			},
			wantErr:     true,
			expectedErr: models.ErrGroupNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			db, teardown := dbInitializedForAccount(t)
			defer teardown()

			model := AccountModel{
				db: db,
			}
			err := model.Create(tt.input.name, tt.input.description, tt.input.saldo, tt.input.groupId)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}
			is.NoErr(err) // got error, did not expect it

			var gotName, gotDescription string
			var gotSaldo float64
			var gotGroupId int

			err = db.QueryRow("SELECT name,description,saldo,group_id FROM accounts WHERE id=?", 1).Scan(
				&gotName, &gotDescription, &gotSaldo, &gotGroupId)
			is.NoErr(err) // got scan error

			is.Equal(gotName, tt.want.name)
			is.Equal(gotDescription, tt.want.description)
			is.Equal(gotSaldo, tt.want.saldo)
			is.Equal(gotGroupId, tt.want.groupId)
		})
	}
}

func TestAccountModel_Read(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	t.Run("read account", func(t *testing.T) {
		is := is.New(t)
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		want := models.Account{
			ID:          1,
			Name:        "tim",
			Description: "",
			Saldo:       12,
			Group: models.Group{
				ID:   1,
				Name: "testgroup1",
			},
		}

		insertTestAccount(t, db, want)

		model := AccountModel{
			db: db,
		}

		got, err := model.Read(1)
		is.NoErr(err) // got error from read, did not expect it

		is.Equal(*got, want)
	})
}

func TestAccountModel_Update(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	tests := []struct {
		name        string
		inital      models.Account
		want        models.Account
		wantErr     bool
		expectedErr error
	}{
		{
			name: "update account",
			inital: models.Account{
				ID:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: models.Group{
					ID:          1,
					Name:        "testgroup1",
					Description: "",
				},
			},
			want: models.Account{
				ID:          1,
				Name:        "tim",
				Description: "descr",
				Saldo:       123,
				Group: models.Group{
					ID:          1,
					Name:        "testgroup1",
					Description: "",
				},
			},
		},
		{
			name: "update account with non existent group",
			inital: models.Account{
				ID:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: models.Group{
					ID:          1,
					Name:        "testgroup1",
					Description: "",
				},
			},
			want: models.Account{
				ID:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: models.Group{
					ID:          12,
					Name:        "testgroup1",
					Description: "",
				},
			},
			wantErr:     true,
			expectedErr: models.ErrGroupNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			db, teardown := dbInitializedForAccount(t)
			defer teardown()

			insertTestAccount(t, db, tt.inital)

			model := AccountModel{
				db: db,
			}

			err := model.Update(&tt.want)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}

			is.NoErr(err) // got error from read, did not expect it

			var gotName, gotDescription string
			var gotSaldo float64
			var gotGroupId int

			err = db.QueryRow("SELECT name,description,saldo,group_id FROM accounts WHERE id=?", 1).Scan(
				&gotName, &gotDescription, &gotSaldo, &gotGroupId)
			is.NoErr(err) // got scan error

			is.Equal(gotName, tt.want.Name)
			is.Equal(gotDescription, tt.want.Description)
			is.Equal(gotSaldo, tt.want.Saldo)
			is.Equal(gotGroupId, tt.want.Group.ID)

		})
	}
}

func TestAccountModel_Delete(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	tests := []struct {
		name         string
		obj          *models.Account
		insertBefore bool
	}{
		{
			name: "delete account",
			obj: &models.Account{
				ID:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: models.Group{
					ID: 1,
				},
			},
			insertBefore: true,
		},
		{
			name: "delete account that does not exist",
			obj: &models.Account{
				ID: 1,
			},
			insertBefore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			db, teardown := dbInitializedForAccount(t)
			defer teardown()

			if tt.insertBefore {
				insertTestAccount(t, db, *tt.obj)
			}

			model := AccountModel{
				db: db,
			}
			err := model.Delete(tt.obj.ID)
			is.NoErr(err)

			var dbName string
			err = db.QueryRow("SELECT name from accounts WHERE id=?", tt.obj.ID).Scan(&dbName)

			if err == nil {
				t.Fatalf("expected err, got none")
			}

			if err != sql.ErrNoRows {
				t.Errorf("got err %v, wanted %v", err, sql.ErrNoRows)
			}
		})
	}
}

func insertTestAccount(t *testing.T, db *sql.DB, account models.Account) {
	t.Helper()

	_, _ = db.Exec("INSERT INTO accounts (id, name, description, saldo, group_id) VALUES (?,?,?,?,?)",
		account.ID,
		account.Name,
		account.Description,
		account.Saldo,
		account.Group.ID,
	)
}

func dbInitializedForAccount(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	db, teardown := getTestDb(t)
	setupScript, _ := ioutil.ReadFile("../testdata/account.sql")
	_, err := db.Exec(string(setupScript))
	if err != nil {
		t.Fatalf("got error initializing account into database: %v", err)
	}
	return db, teardown
}
