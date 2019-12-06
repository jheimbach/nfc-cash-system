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
	ErrSomethingWentWrong    = errors.New("something went wrong")
)

type accountserver struct {
	storage models.AccountStorager
}

func (a *accountserver) List(ctx context.Context, req *api.AccountListRequest) (*api.Accounts, error) {
	accounts, err := a.storage.GetAll()

	if err != nil {
		return nil, ErrGetAll
	}

	return accounts, nil
}

func (a *accountserver) Create(ctx context.Context, req *api.Account) (*api.Account, error) {
	account, err := a.storage.Create(req.Name, req.Description, req.Saldo, req.Group, req.NfcChipId)
	if err != nil {
		return nil, ErrCouldNotCreateAccount
	}

	return account, nil
}

func (a *accountserver) Get(ctx context.Context, req *api.IdRequest) (*api.Account, error) {
	account, err := a.storage.Read(req.Id)

	if err != nil {
		return nil, ErrNotFound
	}

	return account, nil
}

func (a *accountserver) Update(ctx context.Context, req *api.Account) (*api.Account, error) {
	err := a.storage.Update(req)

	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	return req, nil
}

func (a *accountserver) Delete(ctw context.Context, req *api.IdRequest) (*api.Status, error) {
	err := a.storage.Delete(req.Id)

	if err != nil {
		return &api.Status{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &api.Status{Success: true}, nil
}
