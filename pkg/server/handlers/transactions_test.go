package handlers

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jheimbach/nfc-cash-system/api"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test/mock"
	"github.com/jheimbach/nfc-cash-system/pkg/server/repositories"
)

func TestTransactionServer_ListTransactions(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.ListTransactionRequest
		want      *api.ListTransactionsResponse
		wantErr   error
		returnErr error
	}{
		{
			name:  "return all transactions",
			input: &api.ListTransactionRequest{},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(2, 1),
				TotalCount:   2,
			},
		},
		{
			name:      "storage returns error",
			input:     &api.ListTransactionRequest{},
			wantErr:   ErrSomethingWentWrong,
			returnErr: errors.New("test error"),
		},
		{
			name: "return transactions with limit",
			input: &api.ListTransactionRequest{
				Paging: &api.Paging{
					Limit: 5,
				},
			},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(5, 1),
				TotalCount:   5,
			},
		},
		{
			name: "return transactions with limit and offset",
			input: &api.ListTransactionRequest{
				Paging: &api.Paging{
					Limit:  5,
					Offset: 5,
				},
			},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(10, 1)[5:10],
				TotalCount:   5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := transactionServer{
				storage: &mock.TransactionRepository{
					GetAllFunc: func(accountId int32, order string, limit, offset int32) ([]*api.Transaction, int, error) {
						if tt.returnErr != nil {
							return nil, 0, tt.returnErr
						}

						if tt.input.Order != order {
							t.Errorf("got order %q, expected %q", order, tt.input.Order)
						}

						if tt.input.Paging != nil {
							if limit != tt.input.Paging.Limit {
								t.Errorf("got limit %d, expected %d", limit, tt.input.Paging.Limit)
							}
							if offset != tt.input.Paging.Offset {
								t.Errorf("got offset %d, expected %d", offset, tt.input.Paging.Offset)
							}
						}

						return tt.want.Transactions, len(tt.want.Transactions), nil
					},
				},
			}
			got, err := server.ListTransactions(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestTransactionServer_ListTransactionsByAccount(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.ListTransactionsByAccountRequest
		want      *api.ListTransactionsResponse
		wantErr   error
		returnErr error
	}{
		{
			name:  "return all account transaction",
			input: &api.ListTransactionsByAccountRequest{AccountId: 1},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(3, 1),
				TotalCount:   3,
			},
		},
		{
			name: "return account transaction with limit",
			input: &api.ListTransactionsByAccountRequest{
				AccountId: 1,
				Paging: &api.Paging{
					Limit: 5,
				},
			},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(5, 1),
				TotalCount:   5,
			},
		},
		{
			name: "return account transaction with limit and offset",
			input: &api.ListTransactionsByAccountRequest{
				AccountId: 1,
				Paging: &api.Paging{
					Limit:  3,
					Offset: 2,
				},
			},
			want: &api.ListTransactionsResponse{
				Transactions: genTransactionModels(5, 1)[2:5],
				TotalCount:   3,
			},
		},
		{
			name: "return error",
			input: &api.ListTransactionsByAccountRequest{
				AccountId: -45,
			},
			wantErr:   ErrSomethingWentWrong,
			returnErr: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &transactionServer{
				storage: &mock.TransactionRepository{
					GetAllFunc: func(accountId int32, order string, limit, offset int32) ([]*api.Transaction, int, error) {
						if accountId != tt.input.AccountId {
							t.Fatalf("got accountid %d, expected %d", accountId, tt.input.AccountId)
						}
						if tt.returnErr != nil {
							return nil, 0, tt.returnErr
						}

						if tt.input.Order != order {
							t.Errorf("got order %q, expected %q", order, tt.input.Order)
						}
						if tt.input.Paging != nil {
							if limit != tt.input.Paging.Limit {
								t.Errorf("got limit %d, expected %d", limit, tt.input.Paging.Limit)
							}
							if offset != tt.input.Paging.Offset {
								t.Errorf("got offset %d, expected %d", offset, tt.input.Paging.Offset)
							}
						}

						return tt.want.Transactions, len(tt.want.Transactions), nil

					},
				},
			}
			got, err := server.ListTransactionsByAccount(context.Background(), tt.input)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("got err %v, did not expect one", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionServer_CreateTransaction(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.CreateTransactionRequest
		wantErr   error
		returnErr error
	}{
		{
			name: "create transaction",
			input: &api.CreateTransactionRequest{
				Amount:    5,
				AccountId: 1,
			},
		},
		{
			name: "storage returns AccountNotFound",
			input: &api.CreateTransactionRequest{
				Amount:    -5,
				AccountId: 100,
			},
			returnErr: repositories.ErrAccountNotFound,
			wantErr:   ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := transactionServer{
				storage: &mock.TransactionRepository{
					CreateFunc: func(amount float64, accountId int32) (*api.Transaction, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return &api.Transaction{
							Id:       1,
							Amount:   amount,
							OldSaldo: 120,
							NewSaldo: 120 - amount,
							Account:  &api.Account{Id: accountId},
							Created:  timeStamp(),
						}, nil
					},
				},
			}

			got, err := server.CreateTransaction(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}
			want := &api.Transaction{
				Id:       1,
				OldSaldo: 120,
				NewSaldo: 120 - tt.input.Amount,
				Amount:   tt.input.Amount,
				Created:  timeStamp(),
				Account:  &api.Account{Id: 1},
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, expected %v", got, want)
			}
		})
	}
}

func TestTransactionServer_GetTransaction(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.GetTransactionRequest
		want      *api.Transaction
		wantErr   error
		returnErr error
	}{
		{
			name: "returns transaction",
			input: &api.GetTransactionRequest{
				Id:        1,
				AccountId: 1,
			},
			want: genTransactionModels(1, 1)[0],
		},
		{
			name: "returns transaction but account id does not match",
			input: &api.GetTransactionRequest{
				Id:        1,
				AccountId: 1,
			},
			want:    genTransactionModels(1, 2)[0],
			wantErr: ErrTransactionNotFound,
		},
		{
			name: "storage returns error",
			input: &api.GetTransactionRequest{
				Id:        1,
				AccountId: 1,
			},
			want:      genTransactionModels(1, 1)[0],
			wantErr:   ErrSomethingWentWrong,
			returnErr: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := transactionServer{
				storage: &mock.TransactionRepository{
					ReadFunc: func(id int32) (*api.Transaction, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return tt.want, nil
					},
				},
			}

			got, err := server.GetTransaction(context.Background(), tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, expected %v", got, tt.want)
			}
		})
	}
}

func genTransactionModels(num int, accountId int32) []*api.Transaction {
	account := &api.Account{Id: accountId}
	var transactions []*api.Transaction
	for i := 0; i < num; i++ {
		transactions = append(transactions, &api.Transaction{
			Id:       int32(i + 1),
			OldSaldo: 120,
			NewSaldo: 115,
			Amount:   5,
			Created:  timeStamp(),
			Account:  account,
		})
	}

	return transactions
}

func timeStamp() *timestamp.Timestamp {
	created, _ := ptypes.TimestampProto(
		time.Date(2019, 01, 17, 16, 15, 14, 0, time.UTC),
	)
	return created
}
