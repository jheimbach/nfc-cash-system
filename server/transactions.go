package server

import (
	"context"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"google.golang.org/grpc"
)

type transactionServer struct {
	storage models.TransactionStorager
}

func RegisterTransactionServer(server *grpc.Server, storage models.TransactionStorager) {
	api.RegisterTransactionsServiceServer(server, &transactionServer{storage: storage})
}

func (t *transactionServer) All(ctx context.Context, req *api.TransactionAllRequest) (*api.Transactions, error) {
	transactions, err := t.storage.GetAll()
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return transactions, nil
}

func (t *transactionServer) List(ctx context.Context, req *api.TransactionListRequest) (*api.Transactions, error) {
	transactions, err := t.storage.GetAllByAccount(req.AccountId)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return transactions, nil
}

func (t *transactionServer) Create(ctx context.Context, req *api.Transaction) (*api.Transaction, error) {
	transaction, err := t.storage.Create(req.Amount, req.OldSaldo, req.NewSaldo, req.Account)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return transaction, nil
}

func (t *transactionServer) Get(ctx context.Context, req *api.TransactionRequest) (*api.Transaction, error) {
	transaction, err := t.storage.Read(req.Id)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	if transaction.Account.Id != req.AccountId {
		return nil, ErrNotFound
	}

	return transaction, nil
}
