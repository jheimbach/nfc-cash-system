package server

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
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

func Test_accountserver_List(t *testing.T) {
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
