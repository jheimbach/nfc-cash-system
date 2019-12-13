package server

import (
	"context"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transactionServer struct {
	storage models.TransactionStorager
}

func RegisterTransactionServer(server *grpc.Server, storage models.TransactionStorager) {
	api.RegisterTransactionsServiceServer(server, &transactionServer{storage: storage})
}

func (t *transactionServer) ListTransactions(ctx context.Context, req *api.ListTransactionRequest) (*api.ListTransactionsResponse, error) {
	limit, offset := pagingOptions(req.Paging)

	transactions, err := t.storage.GetAll(req.Order, limit, offset)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return &api.ListTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(len(transactions)),
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

	transactions, count, err := t.storage.GetAll(req.AccountId, req.Order, limit, offset)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return &api.ListTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(len(transactions)),
	}, nil
}

func (t *transactionServer) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.Transaction, error) {
	transaction, err := t.storage.Create(req.Amount, req.AccountId)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return transaction, nil
}

func (t *transactionServer) GetTransaction(ctx context.Context, req *api.GetTransactionRequest) (*api.Transaction, error) {
	transaction, err := t.storage.Read(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	if transaction.Account.Id != req.AccountId {
		return nil, status.Error(codes.NotFound, ErrNotFound.Error())
	}

	return transaction, nil
}
