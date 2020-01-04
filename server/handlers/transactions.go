package handlers

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

func (t *transactionServer) ListTransactions(ctx context.Context, req *api.ListTransactionRequest) (*api.ListTransactionsResponse, error) {
	limit, offset := pagingOptions(req.Paging)

	transactions, count, err := t.storage.GetAll(ctx, 0, req.Order, limit, offset)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	return &api.ListTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(count),
	}, nil
}

func pagingOptions(req *api.Paging) (int32, int32) {
	var limit, offset int32

	if req != nil {
		if req.Limit > 0 {
			limit = req.Limit
			if req.Offset > 0 {
				offset = req.Offset
			}
		}
	}
	return limit, offset
}

func (t *transactionServer) ListTransactionsByAccount(ctx context.Context, req *api.ListTransactionsByAccountRequest) (*api.ListTransactionsResponse, error) {
	limit, offset := pagingOptions(req.Paging)

	transactions, count, err := t.storage.GetAll(ctx, req.AccountId, req.Order, limit, offset)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	return &api.ListTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(count),
	}, nil
}

func (t *transactionServer) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.Transaction, error) {
	transaction, err := t.storage.Create(ctx, req.Amount, req.AccountId)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	return transaction, nil
}

func (t *transactionServer) GetTransaction(ctx context.Context, req *api.GetTransactionRequest) (*api.Transaction, error) {
	transaction, err := t.storage.Read(ctx, req.Id)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	if transaction.Account.Id != req.AccountId {
		return nil, ErrAccountNotFound
	}

	return transaction, nil
}
