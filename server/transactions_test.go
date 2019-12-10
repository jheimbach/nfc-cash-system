package server

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type transactionMockStorage struct {
	getAll          func() (*api.Transactions, error)
	read            func(id int32) (*api.Transaction, error)
	create          func(amount, oldSaldo, newSaldo float64, account *api.Account) (*api.Transaction, error)
	getAllByAccount func(accountId int32) (*api.Transactions, error)
}

func (t *transactionMockStorage) GetAll() (*api.Transactions, error) {
	return t.getAll()
}

func (t *transactionMockStorage) Read(id int32) (*api.Transaction, error) {
	return t.read(id)
}

func (t *transactionMockStorage) Create(amount, oldSaldo, newSaldo float64, account *api.Account) (*api.Transaction, error) {
	return t.create(amount, oldSaldo, newSaldo, account)
}

func (t *transactionMockStorage) GetAllByAccount(accountId int32) (*api.Transactions, error) {
	return t.getAllByAccount(accountId)
}

func TestTransactionServer_All(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.TransactionAllRequest
		want      *api.Transactions
		wantErr   error
		returnErr error
	}{
		{
			name:  "return all transactions",
			input: &api.TransactionAllRequest{},
			want: &api.Transactions{
				Transactions: genTransactionModels(2, 1),
			},
		},
		{
			name:      "storage returns error",
			input:     &api.TransactionAllRequest{},
			wantErr:   ErrSomethingWentWrong,
			returnErr: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := transactionServer{
				storage: &transactionMockStorage{
					getAll: func() (*api.Transactions, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return tt.want, nil
					},
				},
			}
			got, err := server.All(context.Background(), tt.input)

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

func TestTransactionServer_List(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.TransactionListRequest
		want      *api.Transactions
		wantErr   error
		returnErr error
	}{
		{
			name:  "return all transaction for the same account",
			input: &api.TransactionListRequest{AccountId: 1},
			want: &api.Transactions{
				Transactions: genTransactionModels(3, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &transactionServer{
				storage: &transactionMockStorage{
					getAllByAccount: func(accountId int32) (*api.Transactions, error) {
						if accountId != tt.input.AccountId {
							t.Fatalf("got accountid %d, expected %d", accountId, tt.input.AccountId)
						}
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return tt.want, nil
					},
				},
			}
			got, err := server.List(context.Background(), tt.input)
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

func TestTransactionServer_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.Transaction
		wantErr   error
		returnErr error
	}{
		{
			name: "create transaction",
			input: &api.Transaction{
				OldSaldo: 120,
				NewSaldo: 115,
				Amount:   -5,
				Account:  &api.Account{Id: 1},
			},
		},
		{
			name: "storage returns AccountNotFound",
			input: &api.Transaction{
				OldSaldo: 120,
				NewSaldo: 115,
				Amount:   -5,
				Account:  &api.Account{Id: 1},
			},
			returnErr: models.ErrAccountNotFound,
			wantErr:   ErrSomethingWentWrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := transactionServer{
				storage: &transactionMockStorage{
					create: func(amount, oldSaldo, newSaldo float64, account *api.Account) (*api.Transaction, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return &api.Transaction{
							Id:       1,
							Amount:   amount,
							OldSaldo: oldSaldo,
							NewSaldo: newSaldo,
							Account:  account,
							Created:  timeStamp(),
						}, nil
					},
				},
			}

			got, err := server.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}
			want := tt.input
			want.Id = 1
			want.Created = timeStamp()

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, expected %v", got, want)
			}
		})
	}
}

func TestTransactionServer_Get(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.TransactionRequest
		want      *api.Transaction
		wantErr   error
		returnErr error
	}{
		{
			name: "returns transaction",
			input: &api.TransactionRequest{
				Id:        1,
				AccountId: 1,
			},
			want: genTransactionModels(1, 1)[0],
		},
		{
			name: "returns transaction but account id does not match",
			input: &api.TransactionRequest{
				Id:        1,
				AccountId: 1,
			},
			want:    genTransactionModels(1, 2)[0],
			wantErr: ErrNotFound,
		},
		{
			name: "storage returns error",
			input: &api.TransactionRequest{
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
				storage: &transactionMockStorage{
					read: func(id int32) (*api.Transaction, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return tt.want, nil
					},
				},
			}

			got, err := server.Get(context.Background(), tt.input)

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
