package models

import (
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
	Create(name, description string, startSaldo float64, group *api.Group, nfcChipId string) (*api.Account, error)
	GetAll() (*api.Accounts, error)
	Read(id int32) (*api.Account, error)
	Delete(id int32) error
	Update(m *api.Account) error
}

type GroupStorager interface {
	Create(name, description string, canOverdraw bool) (*api.Group, error)
	GetAll() (*api.Groups, error)
	Read(id int32) (*api.Group, error)
	Update(group *api.Group) (*api.Group, error)
	Delete(id int32) error
}

type TransactionStorager interface {
	GetAll() (*api.Transactions, error)
	Read(id int32) (*api.Transaction, error)
	Create(amount, oldSaldo, newSaldo float64, account *api.Account) (*api.Transaction, error)
	GetAllByAccount(accountId int32) (*api.Transactions, error)
}
