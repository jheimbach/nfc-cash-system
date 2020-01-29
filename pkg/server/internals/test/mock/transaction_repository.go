package mock

import (
	"context"

	"github.com/jheimbach/nfc-cash-system/api"
)

type TransactionRepository struct {
	CreateFunc             func(float64, int32) (*api.Transaction, error)
	GetAllFunc             func(int32, string, int32, int32) ([]*api.Transaction, int, error)
	ReadFunc               func(int32) (*api.Transaction, error)
	DeleteAllByAccountFunc func(int32) error
}

func (t *TransactionRepository) Create(_ context.Context, amount float64, accountId int32) (*api.Transaction, error) {
	return t.CreateFunc(amount, accountId)
}

func (t *TransactionRepository) GetAll(_ context.Context, accountId int32, order string, limit, offset int32) ([]*api.Transaction, int, error) {
	return t.GetAllFunc(accountId, order, limit, offset)
}

func (t *TransactionRepository) Read(_ context.Context, id int32) (*api.Transaction, error) {
	return t.ReadFunc(id)
}

func (t *TransactionRepository) DeleteAllByAccount(_ context.Context, accountId int32) error {
	return t.DeleteAllByAccountFunc(accountId)
}
