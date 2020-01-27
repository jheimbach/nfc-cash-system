package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test/mock"
	"github.com/JHeimbach/nfc-cash-system/server/repositories"
	isPkg "github.com/matryer/is"
)

var (
	mockGroupOne = &api.Group{
		Id:          1,
		Name:        "testgroup1",
		Description: "",
		CanOverdraw: false,
	}
	mockGroupTwo = &api.Group{
		Id:          2,
		Name:        "testgroup2",
		Description: "with description",
		CanOverdraw: false,
	}
	mockGroupMap = map[int32]*api.Group{
		mockGroupOne.Id: mockGroupOne,
		mockGroupTwo.Id: mockGroupTwo,
	}
)

type groupModelMock struct {
	test   *testing.T
	groups map[int32]*api.Group
}

func (g *groupModelMock) GetAllByIds(_ context.Context, ids []int32) (map[int32]*api.Group, error) {
	m := make(map[int32]*api.Group, len(ids))

	for _, id := range ids {
		if group, ok := g.groups[id]; ok {
			m[id] = group
		}
	}

	return m, nil
}

func (g *groupModelMock) Create(_ context.Context, _, _ string, _ bool) (*api.Group, error) {
	g.test.Fatalf("create of groupmodelmock is not implemented and should not be used")
	return nil, nil
}

func (g *groupModelMock) GetAll(_ context.Context, _, _ int32) ([]*api.Group, int, error) {
	if len(g.groups) < 1 {
		return nil, 0, models.ErrNotFound
	}
	var groups []*api.Group
	for _, group := range g.groups {
		groups = append(groups, group)
	}
	return groups, len(groups), nil
}

func (g *groupModelMock) Read(_ context.Context, id int32) (*api.Group, error) {
	if group, ok := g.groups[id]; ok {
		return group, nil
	}

	return nil, models.ErrGroupNotFound
}

func (g *groupModelMock) Update(_ context.Context, _ *api.Group) (*api.Group, error) {
	g.test.Fatalf("update of groupmodelmock is not implemented and should not be used")
	return nil, nil
}

func (g *groupModelMock) Delete(_ context.Context, _ int32) error {
	g.test.Fatalf("delete of groupmodelmock is not implemented and should not be used")
	return nil
}

func TestAccountModel_Create(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	tests := []struct {
		name          string
		accountCreate *api.CreateAccountRequest
		want          *api.Account
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "creates account",
			accountCreate: &api.CreateAccountRequest{
				Name:        "tim",
				Description: "",
				Saldo:       12,
				GroupId:     1,
				NfcChipId:   "teststringteststring",
			},
			want: &api.Account{
				Id:          1,
				Name:        "tim",
				Description: "",
				Saldo:       12,
				NfcChipId:   "teststringteststring",
				Group:       mockGroupOne,
			},
		},
		{
			name: "creates account but group does not exists",
			accountCreate: &api.CreateAccountRequest{
				Name:        "tim",
				Description: "",
				Saldo:       12,
				GroupId:     100,
				NfcChipId:   "teststring",
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
				groups: &groupModelMock{
					test:   t,
					groups: mockGroupMap,
				},
			}
			account, err := model.Create(context.Background(), tt.accountCreate.Name, tt.accountCreate.Description, tt.accountCreate.Saldo, tt.accountCreate.GroupId, tt.accountCreate.NfcChipId)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}
			is.NoErr(err) // got error, did not expect it

			if tt.want != nil {
				is.Equal(account, tt.want) // returned object is incorrect
			}

			got := &api.Account{Group: mockGroupOne}
			var gotDescription sql.NullString

			err = db.QueryRow("SELECT id,name,description,saldo,nfc_chip_uid FROM accounts WHERE id=?", 1).Scan(
				&got.Id, &got.Name, &gotDescription, &got.Saldo, &got.NfcChipId)
			is.NoErr(err) // got scan error
			got.Description = decodeNullableString(gotDescription)

			is.Equal(got, tt.want)

		})
	}

	t.Run("try to insert new account with same NfcChipId", func(t *testing.T) {
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		model := AccountModel{
			db: db,
			groups: &groupModelMock{
				test:   t,
				groups: mockGroupMap,
			},
		}

		insertTestAccount(t, db, api.Account{
			Id:          1,
			Name:        "tim",
			Description: "",
			Saldo:       12,
			NfcChipId:   "same_id",
			Group:       mockGroupOne,
		})
		_, err := model.Create(context.Background(), "another tim", "", 0, 1, "same_id")
		if err != nil && err != models.ErrDuplicateNfcChipId {
			t.Errorf("got err %q, expected %q", err, models.ErrDuplicateNfcChipId)
		}
	})
}

func TestAccountModel_Read(t *testing.T) {
	test.IsIntegrationTest(t)
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
			Group:       mockGroupOne,
		}

		insertTestAccount(t, db, *want)

		model := AccountModel{
			db: db,
			groups: &groupModelMock{
				test:   t,
				groups: mockGroupMap,
			},
		}

		got, err := model.Read(context.Background(), 1)
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
			Group: mockGroupOne,
		}

		insertTestAccount(t, db, *want)

		model := AccountModel{
			db: db,
			groups: &groupModelMock{
				test:   t,
				groups: mockGroupMap,
			},
		}

		got, err := model.Read(context.Background(), 1)
		is.NoErr(err) // got error from read, did not expect it

		is.Equal(got, want)
	})

	t.Run("read account that does not exist", func(t *testing.T) {
		db, teardown := dbInitializedForAccount(t)
		defer teardown()

		model := AccountModel{
			db: db,
			groups: &groupModelMock{
				test:   t,
				groups: mockGroupMap,
			},
		}

		_, err := model.Read(context.Background(), 100)

		if err != models.ErrNotFound {
			t.Errorf("got %v expected %v", err, models.ErrNotFound)
		}
	})
}

func TestAccountModel_Update(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)

	tests := []struct {
		name        string
		inital      api.Account
		input       api.Account
		want        api.Account
		wantErr     bool
		expectedErr error
	}{
		{
			name: "update description",
			inital: api.Account{
				Id:    1,
				Group: mockGroupOne,
			},
			input: api.Account{
				Id:          1,
				Description: "descr",
				Group: &api.Group{
					Id: 1,
				},
			},
			want: api.Account{
				Id:          1,
				Description: "descr",
				Group:       mockGroupOne,
			},
		},
		{
			name: "update name",
			inital: api.Account{
				Id:    1,
				Name:  "tim",
				Group: mockGroupOne,
			},
			input: api.Account{
				Id:   1,
				Name: "timothy",
				Group: &api.Group{
					Id: 1,
				},
			},
			want: api.Account{
				Id:    1,
				Name:  "timothy",
				Group: mockGroupOne,
			},
		},
		{
			name: "update nfc chip id",
			inital: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupOne,
			},
			input: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip2",
				Group:     mockGroupOne,
			},
			want: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip2",
				Group:     mockGroupOne,
			},
		},
		{
			name: "update saldo 0 is ignored",
			inital: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupOne,
			},
			input: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     0,
				NfcChipId: "testnfcchip",
				Group:     mockGroupOne,
			},
			want: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupOne,
			},
		},
		{
			name: "update group",
			inital: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupOne,
			},
			input: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupTwo,
			},
			want: api.Account{
				Id:        1,
				Name:      "tim",
				Saldo:     123,
				NfcChipId: "testnfcchip",
				Group:     mockGroupTwo,
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
			input: api.Account{
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
		{
			name: "update saldo returns error",
			inital: api.Account{
				Id:    1,
				Name:  "tim",
				Saldo: 12,
				Group: &api.Group{
					Id: 1,
				},
			},
			input: api.Account{
				Id:          1,
				Name:        "tim",
				Description: "",
				Saldo:       120,
				Group:       mockGroupOne,
			},
			wantErr:     true,
			expectedErr: models.ErrUpdateSaldo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			db, teardown := dbInitializedForAccount(t)
			defer teardown()

			insertTestAccount(t, db, tt.inital)

			model := AccountModel{
				db:     db,
				groups: NewGroupModel(db),
			}

			got, err := model.Update(context.Background(), &tt.input)

			if tt.wantErr {
				is.Equal(err, tt.expectedErr) // got not the expected error
				return
			}

			is.NoErr(err) // got error from read, did not expect it

			is.Equal(got, &tt.want) // accounts dont match
		})
	}
}

func TestAccountModel_Delete(t *testing.T) {
	test.IsIntegrationTest(t)
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
			err := model.Delete(context.Background(), tt.obj.Id)
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
	test.IsIntegrationTest(t)
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
			err := model.UpdateSaldo(context.Background(), &tt.obj, tt.newSaldo)
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
	test.IsIntegrationTest(t)

	is := isPkg.New(t)
	db, dbTeardown := dbInitializedForAccountLists(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	type args struct {
		groupId, limit, offset int32
	}
	tests := []struct {
		name      string
		input     args
		want      []*api.Account
		wantErr   error
		wantCount int
	}{
		{
			name:      "return all accounts",
			input:     args{0, 0, 0},
			want:      mockListAccounts(9),
			wantCount: 9,
		},
		{
			name:      "return all accounts with group id 1",
			input:     args{1, 0, 0},
			want:      mockListAccounts(5),
			wantCount: 5,
		},
		{
			name:      "return all accounts limit 5",
			input:     args{0, 5, 0},
			want:      mockListAccounts(5),
			wantCount: 9,
		},
		{
			name:      "return all accounts limit 2 offset 3",
			input:     args{0, 2, 3},
			want:      mockListAccounts(5)[3:5],
			wantCount: 9,
		},
		{
			name:      "return all accounts in group 1 limit 1 offset 2",
			input:     args{1, 1, 2},
			want:      mockListAccounts(5)[2:3],
			wantCount: 5,
		},
	}

	model := AccountModel{
		db: db,
		groups: &groupModelMock{
			test:   t,
			groups: mockGroupMap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			accounts, count, err := model.GetAll(context.Background(), tt.input.groupId, tt.input.limit, tt.input.offset)
			is.NoErr(err)
			is.Equal(accounts, tt.want)
			is.Equal(count, tt.wantCount)
		})
	}
}

func TestAccountModel_GetAllByIds(t *testing.T) {
	test.IsIntegrationTest(t)

	is := isPkg.New(t)
	db, dbTeardown := dbInitializedForAccountLists(t)
	defer func() {
		dbTeardown()
		db.Close()
	}()

	tests := []struct {
		name    string
		input   []int32
		want    map[int32]*api.Account
		wantErr error
	}{
		{
			name:  "return two accounts",
			input: []int32{1, 2},
			want: map[int32]*api.Account{
				1: {
					Id:        1,
					Name:      "testaccount1",
					NfcChipId: "chipid1",
					Group:     mockGroupOne,
				},
				2: {
					Id:        2,
					Name:      "testaccount2",
					NfcChipId: "chipid2",
					Group:     mockGroupOne,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			model := AccountModel{
				db: db,
				groups: &groupModelMock{
					test:   t,
					groups: mockGroupMap,
				},
			}
			got, err := model.GetAllByIds(context.Background(), tt.input)
			if tt.wantErr != nil {
				is.Equal(err, tt.wantErr) // errors don't match
				return
			}

			is.NoErr(err)
			is.Equal(got, tt.want) // return values don't match
		})
	}
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
	setup(dataFor("account"))
	return db, teardown
}

func dbInitializedForAccountLists(t *testing.T) (*sql.DB, func()) {
	db, setup, teardown := getTestDb(t)
	setup(dataFor("account"), dataFor("account_list"))

	return db, teardown
}

func mockListAccounts(num int) []*api.Account {
	accounts := make([]*api.Account, 0, num)
	for i := 1; i <= num; i++ {
		accounts = append(accounts, &api.Account{
			Id:        int32(i),
			Name:      fmt.Sprintf("testaccount%d", i),
			NfcChipId: fmt.Sprintf("chipid%d", i),
			Group: func(i int) *api.Group {
				if i > 5 {
					return mockGroupTwo
				}
				return mockGroupOne
			}(i),
		})
	}
	return accounts
}
