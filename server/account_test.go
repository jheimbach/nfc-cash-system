package server

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	isPkg "github.com/matryer/is"
)

type accountMockStorager struct {
	create func(name, description string, startSaldo float64, group *api.Group, nfcChipId string) (*api.Account, error)
	list   func() (*api.Accounts, error)
	read   func(id int32) (*api.Account, error)
	delete func(id int32) error
	update func(m *api.Account) error
}

func (a accountMockStorager) Create(name, description string, startSaldo float64, group *api.Group, nfcChipId string) (*api.Account, error) {
	return a.create(name, description, startSaldo, group, nfcChipId)
}

func (a accountMockStorager) GetAll() (*api.Accounts, error) {
	return a.list()
}

func (a accountMockStorager) Read(id int32) (*api.Account, error) {
	return a.read(id)
}

func (a accountMockStorager) Delete(id int32) error {
	return a.delete(id)
}

func (a accountMockStorager) Update(m *api.Account) error {
	return a.update(m)
}

func TestAccountserver_List(t *testing.T) {
	type args struct {
		ctx         context.Context
		listRequest *api.AccountListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *api.Accounts
		wantErr error
	}{
		{
			name: "get simple list of accounts",
			args: args{
				ctx:         context.Background(),
				listRequest: &api.AccountListRequest{},
			},
			want: &api.Accounts{
				Accounts: genListModels(2),
			},
		},
		{
			name: "has error",
			args: args{
				ctx:         context.Background(),
				listRequest: &api.AccountListRequest{},
			},
			want:    nil,
			wantErr: ErrGetAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &accountserver{
				storage: accountMockStorager{list: func() (*api.Accounts, error) {
					if tt.wantErr != nil {
						return nil, sql.ErrNoRows
					}
					return &api.Accounts{
						Accounts: genListModels(2),
					}, nil
				},
				},
			}
			got, err := a.List(tt.args.ctx, tt.args.listRequest)

			if err != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountserver_Get(t *testing.T) {
	is := isPkg.New(t)

	want := &api.Account{
		Id:          1,
		Name:        "tim",
		Description: "",
		Saldo:       120,
		NfcChipId:   "asdf",
		Group:       &api.Group{Id: 1},
	}
	server := accountserver{storage: accountMockStorager{
		read: func(id int32) (account *api.Account, err error) {
			want.Id = id
			return want, nil
		},
	}}

	got, err := server.Get(context.Background(), &api.IdRequest{Id: 1})
	is.NoErr(err)
	is.Equal(got, want)
}

func TestAccountserver_Create(t *testing.T) {
	mockStorage := make(map[int32]*api.Account)
	var mockId int32 = 1

	is := isPkg.New(t)

	server := accountserver{
		storage: accountMockStorager{
			create: func(name, description string, startSaldo float64, group *api.Group, nfcChipId string) (account *api.Account, err error) {
				acc := &api.Account{
					Id:          mockId,
					Name:        name,
					Description: description,
					Saldo:       startSaldo,
					NfcChipId:   nfcChipId,
					Group:       group,
				}

				mockStorage[mockId] = acc
				mockId = mockId + 1

				return acc, nil
			},
		},
	}

	want := &api.Account{
		Name:        "test",
		Description: "",
		Saldo:       120,
		NfcChipId:   "nfcchip",
		Group: &api.Group{
			Id: 1,
		},
	}

	got, err := server.Create(context.Background(), want)
	is.NoErr(err)

	want.Id = mockId - 1

	is.Equal(got, want)                  // got wrong account back
	is.Equal(want, mockStorage[want.Id]) // account was not created
}

func TestAccountserver_Update(t *testing.T) {
	mockStorage := genMapModels(3)
	is := isPkg.New(t)

	server := accountserver{
		storage: accountMockStorager{
			update: func(m *api.Account) error {
				mockStorage[m.Id] = m
				return nil
			},
		},
	}

	updateAccount := &api.Account{
		Id:          1,
		Name:        "test",
		Description: "test",
		Saldo:       145,
		NfcChipId:   "nfc_chip_1",
		Group: &api.Group{
			Id: 1,
		},
	}

	got, err := server.Update(context.Background(), updateAccount)
	is.NoErr(err)
	is.Equal(got, updateAccount)       // returned account is not the same
	is.Equal(got, mockStorage[got.Id]) // account was not updated
}

func TestAccountserver_Delete(t *testing.T) {
	mockStorage := genMapModels(3)
	is := isPkg.New(t)

	server := accountserver{
		storage: accountMockStorager{
			delete: func(id int32) error {
				delete(mockStorage, id)
				return nil
			},
		},
	}

	want := &api.Status{
		Success:      true,
		ErrorMessage: "",
	}

	got, err := server.Delete(context.Background(), &api.IdRequest{Id: 1})
	is.NoErr(err)
	is.Equal(got, want) // status is not correct
}

func genListModels(num int) []*api.Account {
	accounts := make([]*api.Account, 0, num)

	for i := 1; i <= num; i++ {
		accounts = append(accounts, &api.Account{
			Id:          int32(i),
			Name:        "test",
			Description: "test",
			Saldo:       0,
			NfcChipId:   fmt.Sprintf("ncf_chip_%d", i),
			Group: &api.Group{
				Id: 1,
			},
		})
	}
	return accounts
}

func genMapModels(num int) map[int32]*api.Account {
	accounts := genListModels(num)
	m := make(map[int32]*api.Account)
	for _, account := range accounts {
		m[account.Id] = account
	}
	return m
}
