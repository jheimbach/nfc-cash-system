package models

import (
	"errors"
	"time"

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

type User struct {
	ID      int
	Name    string
	Email   string
	Created time.Time
}

type Transaction struct {
	ID       int
	OldSaldo float64
	NewSaldo float64
	Amount   float64
	Created  time.Time
	Account  api.Account
}

type TransactionPaging struct {
	CurrentPage  int
	MaxPage      int
	Transactions []Transaction
}

type AccountStorager interface {
	Create(name, description string, startSaldo float64, groupId int, nfcChipId string) (*api.Account, error)
	GetAll() ([]*api.Account, error)
	Read(id int) (*api.Account, error)
	Delete(id int) error
	Update(m *api.Account) error
}
