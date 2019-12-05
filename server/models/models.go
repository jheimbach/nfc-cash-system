package models

import (
	"errors"

	"github.com/JHeimbach/nfc-cash-system/server/api"
)

var (
	ErrDuplicateEmail     = errors.New("models: duplicate user email")
	ErrDuplicateNfcChipId = errors.New("models: duplicate account nfc chip id")
	ErrNotFound           = errors.New("models: not found")
	ErrInvalidCredentials = errors.New("models: email or password incorrect")
	ErrModelNotSaved      = errors.New("models: got no id on update, did you mean to create the group")
	ErrNonEmptyDelete     = errors.New("models: can not delete, item is still referenced")
	ErrGroupNotFound      = errors.New("models: group for given id does not exist")
	ErrAccountNotFound    = errors.New("models: account for given id does not exist")
)

type AccountStorager interface {
	Create(name, description string, startSaldo float64, group *api.Group, nfcChipId string) (*api.Account, error)
	GetAll() (*api.Accounts, error)
	Read(id int32) (*api.Account, error)
	Delete(id int32) error
	Update(m *api.Account) error
}
