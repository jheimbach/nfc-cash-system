package models

import (
	"context"
	"errors"

	"github.com/JHeimbach/nfc-cash-system/server/api"
)

var (
	ErrDuplicateEmail      = errors.New("duplicate user email")
	ErrDuplicateNfcChipId  = errors.New("duplicate account nfc chip id")
	ErrNotFound            = errors.New("not found")
	ErrInvalidCredentials  = errors.New("email or password incorrect")
	ErrModelNotSaved       = errors.New("got no id on update, did you mean to create the group")
	ErrNonEmptyDelete      = errors.New("can not delete, item is still referenced")
	ErrGroupNotFound       = errors.New("group for given id does not exist")
	ErrAccountNotFound     = errors.New("account for given id does not exist")
	ErrTransactionNotFound = errors.New("transaction for given id does not exist")
)

type AccountStorager interface {
	Create(ctx context.Context, name, description string, startSaldo float64, groupId int32, nfcChipId string) (*api.Account, error)

	GetAll(ctx context.Context, groupId, limit, offset int32) ([]*api.Account, int, error)
	GetAllByIds(ctx context.Context, ids []int32) (map[int32]*api.Account, error)

	Read(ctx context.Context, id int32) (*api.Account, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, m *api.Account) error

	UpdateSaldo(ctx context.Context, m *api.Account, newSaldo float64) error
}

type GroupStorager interface {
	Create(ctx context.Context, name, description string, canOverdraw bool) (*api.Group, error)

	GetAll(ctx context.Context, limit, offset int32) ([]*api.Group, int, error)
	GetAllByIds(ctx context.Context, ids []int32) (map[int32]*api.Group, error)

	Read(ctx context.Context, id int32) (*api.Group, error)
	Update(ctx context.Context, group *api.Group) (*api.Group, error)
	Delete(ctx context.Context, id int32) error
}

type TransactionStorager interface {
	Create(amount float64, accountId int32) (*api.Transaction, error)

	GetAll(accountId int32, order string, limit, offset int32) ([]*api.Transaction, int, error)

	Read(id int32) (*api.Transaction, error)
}
