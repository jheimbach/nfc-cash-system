package models

import (
	"errors"
	"time"
)

var ErrDuplicateEmail = errors.New("models: duplicate user email")
var ErrNotFound = errors.New("models: not found")
var ErrInvalidCredentials = errors.New("models: email or password incorrect")
var ErrModelNotSaved = errors.New("models: got no id on update, did you mean to create the group")
var ErrNonEmptyDelete = errors.New("models: can not delete, item is still referenced")

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
	Saldo       int
	Group       Group
}
