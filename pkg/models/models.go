package models

import (
	"errors"
	"time"
)

var ErrDuplicateEmail = errors.New("models: duplicate user email")

type User struct {
	ID      int
	Name    string
	Email   string
	Created time.Time
}
