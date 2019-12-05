package mysql

import (
	"database/sql"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
)

func TestAccountModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	type fields struct {
		name, description string
		saldo             float64
		groupId           int
		nfcChipId         string
	}
	tests := []struct {
		name          string
		accountFields fields
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "creates account",
			accountFields: fields{
				name:        "tim",
				description: "",
				saldo:       12,
				groupId:     1,
				nfcChipId:   "teststringteststring",
			},
		},
		{
			name: "creates account but group does not exists",
			accountFields: fields{
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
			_, err := model.Create(tt.accountFields.name, tt.accountFields.description, tt.accountFields.saldo, tt.accountFields.groupId, tt.accountFields.nfcChipId)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}
			is.NoErr(err) // got error, did not expect it

			var gotName, gotNfcChipId string
			var gotDescription sql.NullString
			var gotSaldo float64
			var gotGroupId int

			err = db.QueryRow("SELECT name,description,saldo,group_id,nfc_chip_uid FROM accounts WHERE id=?", 1).Scan(
				&gotName, &gotDescription, &gotSaldo, &gotGroupId, &gotNfcChipId)
			is.NoErr(err) // got scan error

			is.Equal(gotName, tt.accountFields.name)                                     // name does not match
			is.Equal(decodeNullableString(gotDescription), tt.accountFields.description) // Description does not match
			is.Equal(gotSaldo, tt.accountFields.saldo)                                   // Saldo does not match
			is.Equal(gotGroupId, tt.accountFields.groupId)                               // GroupId does not match
			is.Equal(gotNfcChipId, tt.accountFields.nfcChipId)                           // NfcChipId does not match
		})
	}

	t.Run("try to insert new account with same NfcChipId", func(t *testing.T) {
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		model := AccountModel{
			db: db,
		}

		insertTestAccount(t, db, api.Account{
			Id:          1,
			Name:        "tim",
			Description: "",
			Saldo:       12,
			NfcChipId:   "same_id",
			Group: &api.Group{
				Id: 1,
			},
		})
		_, err := model.Create("another tim", "", 0, 1, "same_id")
		if err != nil && err != models.ErrDuplicateNfcChipId {
			t.Errorf("got err %q, expected %q", err, models.ErrDuplicateNfcChipId)
		}
	})
}

func TestAccountModel_Read(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	t.Run("read account", func(t *testing.T) {
		is := is.New(t)
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		want := &api.Account{
			Id:          1,
			Name:        "tim",
			Description: "",
			Saldo:       12,
			NfcChipId:   "testchipid",
			Group: &api.Group{
				Id: 1,
			},
		}

		insertTestAccount(t, db, *want)

		model := AccountModel{
			db: db,
		}

		got, err := model.Read(1)
		is.NoErr(err) // got error from read, did not expect it

		is.Equal(got, want)
	})

	t.Run("read account with null description", func(t *testing.T) {
		is := is.New(t)
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		want := &api.Account{
			Id:    1,
			Name:  "tim",
			Saldo: 12,
			Group: &api.Group{
				Id: 1,
			},
		}

		insertTestAccount(t, db, *want)

		model := AccountModel{
			db: db,
		}

		got, err := model.Read(1)
		is.NoErr(err) // got error from read, did not expect it

		is.Equal(got, want)
	})

	t.Run("read account that does not exist", func(t *testing.T) {
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		model := AccountModel{
			db: db,
		}

		_, err := model.Read(100)

		if err != models.ErrNotFound {
			t.Errorf("got %v expected %v", err, models.ErrNotFound)
		}
	})
}

func TestAccountModel_Update(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	tests := []struct {
		name        string
		inital      api.Account
		want        api.Account
		wantErr     bool
		expectedErr error
	}{
		{
			name: "update account",
			inital: api.Account{
				Id:    1,
				Name:  "tim",
				Saldo: 12,
				Group: &api.Group{
					Id: 1,
				},
			},
			want: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "descr",
				Saldo:       123,
				Group: &api.Group{
					Id: 1,
				},
			},
		},
		{
			name: "update nfc chip id",
			inital: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     12,
				NfcChipId: "testnfcchip",
				Group: &api.Group{
					Id: 1,
				},
			},
			want: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "descr",
				Saldo:       123,
				NfcChipId:   "testnfcchip2",
				Group: &api.Group{
					Id: 1,
				},
			},
		},
		{
			name: "update account with non existent group",
			inital: api.Account{
				Id:    1,
				Name:  "tim",
				Saldo: 12,
				Group: &api.Group{
					Id: 1,
				},
			},
			want: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: &api.Group{
					Id: 12,
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

			err := model.Update(tt.want)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}

			is.NoErr(err) // got error from read, did not expect it

			var got = api.Account{Group: &api.Group{}}
			var nullDescription sql.NullString
			err = db.QueryRow("SELECT id,name,description,saldo,group_id,nfc_chip_uid FROM accounts WHERE id=?", 1).Scan(
				&got.Id, &got.Name, &nullDescription, &got.Saldo, &got.Group.Id, &got.NfcChipId)
			is.NoErr(err) // got scan error

			got.Description = decodeNullableString(nullDescription)

			is.Equal(got.Name, tt.want.Name)               // name does not match
			is.Equal(got.Description, tt.want.Description) // description does not match
			is.Equal(got.Saldo, tt.want.Saldo)             // saldo does not match
			is.Equal(got.Group.Id, tt.want.Group.Id)       // groupId does not match
			is.Equal(got.NfcChipId, tt.want.NfcChipId)     // nfcChipId does not match
		})
	}
}

func TestAccountModel_Delete(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	tests := []struct {
		name         string
		obj          api.Account
		insertBefore bool
	}{
		{
			name: "delete account",
			obj: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				Group: &api.Group{
					Id: 1,
				},
			},
			insertBefore: true,
		},
		{
			name: "delete account that does not exist",
			obj: api.Account{
				Id: 1,
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
				insertTestAccount(t, db, tt.obj)
			}

			model := AccountModel{
				db: db,
			}
			err := model.Delete(tt.obj.Id)
			is.NoErr(err)

			var dbName string
			err = db.QueryRow("SELECT name from accounts WHERE id=?", tt.obj.Id).Scan(&dbName)

			if err == nil {
				t.Fatalf("expected err, got none")
			}

			if err != sql.ErrNoRows {
				t.Errorf("got err %v, wanted %v", err, sql.ErrNoRows)
			}
		})
	}
}

func TestAccountModel_UpdateSaldo(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	tests := []struct {
		name           string
		obj            api.Account
		insertObj      bool
		newSaldo       float64
		expectDbChange bool
	}{
		{
			name: "update saldo",
			obj: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "",
				Saldo:       50,
				Group: &api.Group{
					Id: 1,
				},
			},
			insertObj:      true,
			newSaldo:       65,
			expectDbChange: true,
		},
		{
			name: "update saldo on undefined account",
			obj: api.Account{
				Id: 10,
			},
			newSaldo: 65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			db, teardown := dbInitializedForAccount(t)
			defer teardown()

			if tt.insertObj {
				insertTestAccount(t, db, tt.obj)
			}

			model := AccountModel{db: db}
			err := model.UpdateSaldo(tt.obj.Id, tt.newSaldo)
			is.NoErr(err)

			if tt.expectDbChange {
				var dbSaldo float64
				err = db.QueryRow("SELECT saldo from accounts WHERE id=?", tt.obj.Id).Scan(&dbSaldo)
				is.NoErr(err)
				is.Equal(dbSaldo, tt.newSaldo)
			}
		})
	}
}

func TestAccountModel_GetAll(t *testing.T) {
	isIntegrationTest(t)

	is := isPkg.New(t)
	db, dbTeardown := dbInitializedForAccountLists(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := AccountModel{
		db: db,
	}

	accounts, err := model.GetAll()
	is.NoErr(err)
	is.Equal(len(accounts), 9) // expect 9 accounts
}

func TestAccountModel_GetAllByGroup(t *testing.T) {
	isIntegrationTest(t)

	is := isPkg.New(t)
	db, dbTeardown := dbInitializedForAccountLists(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := AccountModel{
		db: db,
	}

	accounts, err := model.GetAllByGroup(1)
	is.NoErr(err)
	is.Equal(len(accounts), 5)
}

func insertTestAccount(t *testing.T, db *sql.DB, account api.Account) {
	t.Helper()

	_, _ = db.Exec("INSERT INTO accounts (id, name, description, saldo, group_id, nfc_chip_uid) VALUES (?,?,?,?,?,?)",
		account.Id,
		account.Name,
		createNullableString(account.Description),
		account.Saldo,
		account.Group.Id,
		account.NfcChipId,
	)
}

func dbInitializedForAccount(t *testing.T) (*sql.DB, func()) {
	db, setup, teardown := getTestDb(t)
	setup("../testdata/account.sql")
	return db, teardown
}

func dbInitializedForAccountLists(t *testing.T) (*sql.DB, func()) {
	db, setup, teardown := getTestDb(t)
	setup("../testdata/account.sql", "../testdata/account_lists.sql")

	return db, teardown
}
