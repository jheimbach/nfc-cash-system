package mock

import (
	"context"

	"github.com/JHeimbach/nfc-cash-system/server/api"
)

type AccountRepository struct {
	CreateFunc      func(string, string, float64, int32, string) (*api.Account, error)
	GetAllFunc      func(int32, int32, int32) ([]*api.Account, int, error)
	GetAllByIdsFunc func([]int32) (map[int32]*api.Account, error)
	ReadFunc        func(int32) (*api.Account, error)
	DeleteFunc      func(int32) error
	UpdateFunc      func(*api.Account) (*api.Account, error)
	UpdateSaldoFunc func(*api.Account, float64) error
}

func (a *AccountRepository) Create(_ context.Context, name, description string, startSaldo float64, groupId int32, nfcChipId string) (*api.Account, error) {
	return a.CreateFunc(name, description, startSaldo, groupId, nfcChipId)
}

func (a *AccountRepository) GetAll(_ context.Context, groupId, limit, offset int32) ([]*api.Account, int, error) {
	return a.GetAllFunc(groupId, limit, offset)
}

func (a *AccountRepository) GetAllByIds(_ context.Context, ids []int32) (map[int32]*api.Account, error) {
	return a.GetAllByIdsFunc(ids)
}

func (a *AccountRepository) Read(_ context.Context, id int32) (*api.Account, error) {
	return a.ReadFunc(id)
}

func (a *AccountRepository) Delete(_ context.Context, id int32) error {
	return a.DeleteFunc(id)
}

func (a *AccountRepository) Update(_ context.Context, m *api.Account) (*api.Account, error) {
	return a.UpdateFunc(m)
}

func (a *AccountRepository) UpdateSaldo(_ context.Context, m *api.Account, newSaldo float64) error {
	return a.UpdateSaldoFunc(m, newSaldo)
}
