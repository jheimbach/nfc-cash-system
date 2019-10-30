package models

import (
	"errors"
	"time"
)

var ErrDuplicateEmail = errors.New("models: duplicate user email")
var ErrNotFound = errors.New("models: not found")
var ErrInvalidCredentials = errors.New("models: email or password incorrect")

type User struct {
	ID      int
	Name    string
	Email   string
	Created time.Time
}
