package models

import (
	"errors"
	"time"
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

type Group struct {
	ID          int
	Name        string
	Description string
	CanOverDraw bool
}

type Account struct {
	ID          int
	Name        string
	Description string
	Saldo       float64
	NfcChipId   string
	GroupId     int
}

type AccountPaging struct {
	CurrentPage int
	MaxPage     int
	Accounts    []Account
}

type Transaction struct {
	ID       int
	OldSaldo float64
	NewSaldo float64
	Amount   float64
	Created  time.Time
	Account  Account
}

type TransactionPaging struct {
	CurrentPage  int
	MaxPage      int
	Transactions []Transaction
}

type AccountStorager interface {
	Create(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error)
	GetAll() ([]Account, error)
}
