package server

import (
	"context"
	"errors"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func RegisterAccountServer(s *grpc.Server, storage models.AccountStorager) {
	api.RegisterAccountServiceServer(s, &accountserver{storage: storage})
}

func (a *accountserver) ListAccounts(ctx context.Context, req *api.ListAccountsRequest) (*api.ListAccountsResponse, error) {
	var limit int32 = 0
	var offset int32 = 0

	if req.Paging != nil {
		limit = req.Paging.Limit
		offset = req.Paging.Offset
	}
	accounts, totalCount, err := a.storage.GetAll(ctx, req.GroupId, limit, offset)

	if err != nil {
		return nil, ErrGetAll
	}

	return &api.ListAccountsResponse{
		Accounts:   accounts,
		TotalCount: int32(totalCount),
	}, nil
}

func (a *accountserver) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.Account, error) {
	account, err := a.storage.Create(ctx, req.Name, req.Description, req.Saldo, req.GroupId, req.NfcChipId)
	if err != nil {
		return nil, ErrCouldNotCreateAccount
	}

	return account, nil
}

func (a *accountserver) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.Account, error) {
	account, err := a.storage.Read(ctx, req.Id)

	if err != nil {
		return nil, ErrNotFound
	}

	return account, nil
}

func (a *accountserver) UpdateAccount(ctx context.Context, req *api.Account) (*api.Account, error) {
	err := a.storage.Update(ctx, req)

	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	return req, nil
}

func (a *accountserver) DeleteAccount(ctx context.Context, req *api.DeleteAccountRequest) (*empty.Empty, error) {
	err := a.storage.Delete(ctx, req.Id)

	if err != nil {
		if err == models.ErrNotFound {
			return &empty.Empty{}, status.Error(codes.NotFound, ErrNotFound.Error())
		}

		return &empty.Empty{}, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	return &empty.Empty{}, nil
}
