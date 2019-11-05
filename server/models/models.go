package models

import (
	"errors"
	"time"
)

var (
	ErrDuplicateEmail     = errors.New("models: duplicate user email")
	ErrNotFound           = errors.New("models: not found")
	ErrInvalidCredentials = errors.New("models: email or password incorrect")
	ErrModelNotSaved      = errors.New("models: got no id on update, did you mean to create the group")
	ErrNonEmptyDelete     = errors.New("models: can not delete, item is still referenced")
	ErrGroupNotFound      = errors.New("models: group for given group id does not exist")
)

const (
	DefaultPageSize = -1
	DefaultPage     = 0
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
}

type Account struct {
	ID          int
	Name        string
	Description string
	Saldo       float64
	Group       Group
}

type AccountPaging struct {
	CurrentPage int
	MaxPage     int
	Accounts    []*Account
}
