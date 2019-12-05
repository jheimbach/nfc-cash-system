package server

import (
	"context"
	"errors"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
)

var (
	ErrGetAll                = errors.New("could not load list of accounts")
	ErrCouldNotCreateAccount = errors.New("could not save new account")
	ErrNotFound              = errors.New("could not find account")
	ErrCouldNotParseId       = errors.New("could not parse account id")
	ErrSomethingWentWrong    = errors.New("something went wrong")
)

type accountserver struct {
	storage models.AccountStorager
}

func (a *accountserver) List(context.Context, *api.AccountListRequest) (*api.Accounts, error) {
	accounts, err := a.storage.GetAll()

	if err != nil {
		return nil, ErrGetAll
	}

	return accounts, nil
}

func (a *accountserver) Create(context.Context, *api.Account) (*api.Account, error) {
	panic("implement me")
}

func (a *accountserver) Get(context.Context, *api.IdRequest) (*api.Account, error) {
	panic("implement me")
}

func (a *accountserver) Update(context.Context, *api.Account) (*api.Account, error) {
	panic("implement me")
}

func (a *accountserver) Delete(context.Context, *api.IdRequest) (*api.Status, error) {
	panic("implement me")
}
