package mysql

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test/mock"
	"github.com/JHeimbach/nfc-cash-system/server/repositories"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	isPkg "github.com/matryer/is"
)

var (
	accountOne = &api.Account{
		Id:        1,
		Name:      "testaccount",
		Saldo:     12,
		NfcChipId: "testchipid",
		Group: &api.Group{
			Id:   1,
			Name: "testgroup1",
		},
	}
	accountTwo = &api.Account{
		Id:        2,
		Name:      "testaccount1",
		Saldo:     120,
		NfcChipId: "testchipid2",
		Group: &api.Group{
			Id:   1,
			Name: "testgroup1",
		},
	}
	accountMap = map[int32]*api.Account{
		1: accountOne,
		2: accountTwo,
	}
)
var accountMock = &mock.AccountRepository{
	ReadFunc: func(i int32) (account *api.Account, err error) {
		return accountOne, nil
	},

	UpdateSaldoFunc: func(account *api.Account, f float64) error {
		return nil
	},
	GetAllByIdsFunc: func(int32s []int32) (m map[int32]*api.Account, err error) {
		return accountMap, nil
	},
}

func TestTransactionModel_Create(t *testing.T) {
	is, teardown := initTransactionIntegrationTest(t)
	defer teardown()

	tests := []struct {
		name        string
		input       *api.CreateTransactionRequest
		want        *api.Transaction
		wantErr     bool
		expectedErr error
	}{
		{
			name: "create new transaction",
			input: &api.CreateTransactionRequest{
				Amount:    6,
				AccountId: 1,
			},
			want: &api.Transaction{
				Id:       1,
				OldSaldo: 12,
				NewSaldo: 6,
				Amount:   6,
				Account: &api.Account{
					Id:        1,
					Name:      "testaccount",
					Saldo:     6,
					NfcChipId: "testchipid",
					Group: &api.Group{
						Id:   1,
						Name: "testgroup1",
					},
				},
			},
		}, {
			name: "create new transaction with nonexistent account",
			input: &api.CreateTransactionRequest{
				Amount:    6,
				AccountId: 100,
			},
			wantErr:     true,
			expectedErr: repositories.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			td := initDbForTransactions(t)
			defer td()

			got, err := _transactionModel.Create(context.Background(), tt.input.Amount, tt.input.AccountId)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("dbTransaction err %v, expected %v", err, tt.expectedErr)
				}
				return
			}

			is.NoErr(err)

			is.True(got.Created.Nanos != 0) // got.created not set
			tt.want.Created = got.Created   // set timestamp to want for not messing with equality of timestamps

			is.Equal(got, tt.want)

			dbTransaction := api.Transaction{Account: &api.Account{}}
			var created time.Time

			stmt := `SELECT id, new_saldo, old_saldo, amount,created, account_id from transactions WHERE id=?`
			err = _conn.QueryRow(stmt, 1).Scan(
				&dbTransaction.Id, &dbTransaction.NewSaldo, &dbTransaction.OldSaldo, &dbTransaction.Amount, &created, &dbTransaction.Account.Id,
			)
			is.NoErr(err)
			dbTransaction.Created, _ = ptypes.TimestampProto(created)

			is.Equal(dbTransaction.Id, tt.want.Id)                 // id does not match
			is.Equal(dbTransaction.OldSaldo, tt.want.OldSaldo)     // oldSaldo does not match
			is.Equal(dbTransaction.NewSaldo, tt.want.NewSaldo)     // newSaldo does not match
			is.True(!created.IsZero())                             // created is zero, should be timestamp
			is.Equal(dbTransaction.Account.Id, tt.want.Account.Id) // account does not match

		})
	}
}

func TestTransactionModel_Get(t *testing.T) {
	is, teardown := initTransactionIntegrationTest(t)
	defer teardown()

	td := initDbForTransactionList(t)
	defer td()

	t.Run("get a transaction", func(t *testing.T) {
		is := is.New(t)

		_transactionModel.accounts = &mock.AccountRepository{
			ReadFunc: func(id int32) (account *api.Account, err error) {
				return &api.Account{
					Id: id,
				}, nil
			},
		}
		created, _ := ptypes.TimestampProto(time.Date(2019, 01, 17, 16, 15, 14, 0, time.UTC))

		want := &api.Transaction{
			Id:       1,
			OldSaldo: 120,
			NewSaldo: 115,
			Amount:   5,
			Created:  created,
			Account: &api.Account{
				Id: 1,
			},
		}

		transaction, err := _transactionModel.Read(context.Background(), 1)
		is.NoErr(err)

		is.Equal(transaction, want)
	})

	t.Run("no transactions found", func(t *testing.T) {
		is := is.New(t)

		transaction, err := _transactionModel.Read(context.Background(), -45)
		is.Equal(transaction, nil)

		if err != repositories.ErrNotFound {
			t.Errorf("got err %v, expected %v", err, repositories.ErrNotFound)
		}
	})

}

func TestTransactionModel_GetAll(t *testing.T) {
	is, teardown := initTransactionIntegrationTest(t)
	defer teardown()

	type args struct {
		accountId, limit, offset int32
		order                    string
	}
	tests := []struct {
		name      string
		input     args
		want      []*api.Transaction
		wantCount int
	}{
		{
			name:      "get all transactions",
			input:     args{},
			want:      transisitonList(0),
			wantCount: 9,
		},
		{
			name: "get all transactions for account id 1",
			input: args{
				accountId: 1,
			},
			want:      transisitonList(1),
			wantCount: 5,
		},
		{
			name: "no transactions found",
			input: args{
				accountId: 9999,
			},
			wantCount: 0,
		},
		{
			name: "get transactions with limit",
			input: args{
				limit: 5,
			},
			want:      transisitonList(0)[:5],
			wantCount: 9,
		},
		{
			name: "get all transactions for account id 1 with limit",
			input: args{
				accountId: 1,
				limit:     3,
			},
			want:      transisitonList(1)[:3],
			wantCount: 5,
		},
		{
			name: "get transactions with limit and offset",
			input: args{
				limit:  3,
				offset: 2,
			},
			want:      transisitonList(0)[2:5],
			wantCount: 9,
		},
		{
			name: "get all transactions with order DESC",
			input: args{
				order: "DESC",
			},
			want:      transisitonList(0),
			wantCount: 9,
		},
		{
			name: "get all transactions with order desc",
			input: args{
				order: "desc",
			},
			want:      transisitonList(0),
			wantCount: 9,
		},
		{
			name: "get all transactions default order is DESC",
			input: args{
				order: "something invalid",
			},
			want:      transisitonList(0),
			wantCount: 9,
		},
		{
			name: "get all transactions with order ASC",
			input: args{
				order: "ASC",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(0))
				sort.Sort(s)
				return s
			}(),
			wantCount: 9,
		},
		{
			name: "get all transactions with order asc",
			input: args{
				order: "asc",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(0))
				sort.Sort(s)
				return s
			}(),
			wantCount: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			td := initDbForTransactionList(t)
			defer td()

			got, count, err := _transactionModel.GetAll(context.Background(), tt.input.accountId, tt.input.order, tt.input.limit, tt.input.offset)
			is.NoErr(err)
			is.Equal(got, tt.want)
			is.Equal(count, tt.wantCount)
		})
	}

}

func TestTransactionModel_DeleteAllByAccount(t *testing.T) {
	is, teardown := initTransactionIntegrationTest(t)
	defer teardown()

	tests := []struct {
		name      string
		accountId int32
	}{
		{
			name:      "delete all transactions by account id 1",
			accountId: 1,
		},
		{
			name:      "delete all transactions by non existent account",
			accountId: -45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			td := initDbForTransactionList(t)
			defer td()

			err := _transactionModel.DeleteAllByAccount(context.Background(), tt.accountId)
			is.NoErr(err)

			stmt := `SELECT id from transactions WHERE account_id=?`
			rows, err := _conn.Query(stmt, tt.accountId)
			is.NoErr(err) // could not query transactions
			defer rows.Close()
			if rows.Next() {
				t.Errorf("expected empty result")
			}
		})
	}

}

func transisitonList(accountId int) []*api.Transaction {
	transactionsTwo := []*api.Transaction{
		{
			Id:       9,
			OldSaldo: 105,
			NewSaldo: 110,
			Amount:   -5,
			Created:  timeStampMock(9),
			Account:  accountTwo,
		},
		{
			Id:       8,
			OldSaldo: 110,
			NewSaldo: 105,
			Amount:   5,
			Created:  timeStampMock(8),
			Account:  accountTwo,
		},
		{
			Id:       7,
			OldSaldo: 115,
			NewSaldo: 110,
			Amount:   5,
			Created:  timeStampMock(7),
			Account:  accountTwo,
		},
		{
			Id:       6,
			OldSaldo: 120,
			NewSaldo: 115,
			Amount:   5,
			Created:  timeStampMock(6),
			Account:  accountTwo,
		},
	}
	transactionsOne := []*api.Transaction{
		{
			Id:       5,
			OldSaldo: 100,
			NewSaldo: 105,
			Amount:   -5,
			Created:  timeStampMock(5),
			Account:  accountOne,
		},
		{
			Id:       4,
			OldSaldo: 105,
			NewSaldo: 100,
			Amount:   5,
			Created:  timeStampMock(4),
			Account:  accountOne,
		},
		{
			Id:       3,
			OldSaldo: 110,
			NewSaldo: 105,
			Amount:   5,
			Created:  timeStampMock(3),
			Account:  accountOne,
		},
		{
			Id:       2,
			OldSaldo: 115,
			NewSaldo: 110,
			Amount:   5,
			Created:  timeStampMock(2),
			Account:  accountOne,
		},
		{
			Id:       1,
			OldSaldo: 120,
			NewSaldo: 115,
			Amount:   5,
			Created:  timeStampMock(1),
			Account:  accountOne,
		},
	}
	switch accountId {
	case 1:
		return transactionsOne
	case 2:
		return transactionsTwo
	default:
		return append(transactionsTwo, transactionsOne...)
	}
}

type SortTransactions []*api.Transaction

func (s SortTransactions) Less(i, j int) bool {
	timeI, _ := ptypes.Timestamp(s[i].Created)
	timeJ, _ := ptypes.Timestamp(s[j].Created)

	return timeI.Before(timeJ)
}

func (s SortTransactions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortTransactions) Len() int {
	return len(s)
}

func initDbForTransactions(t *testing.T) func() error {
	t.Helper()
	err := test.SetupDB(_conn, dataFor("transaction"))
	if err != nil {
		t.Fatal(err)
	}
	return teardownDB(_conn)
}

func initDbForTransactionList(t *testing.T) func() error {
	t.Helper()
	err := test.SetupDB(_conn, dataFor("transaction"), dataFor("transaction_list"))
	if err != nil {
		t.Fatal(err)
	}
	return teardownDB(_conn)
}

func initTransactionIntegrationTest(t *testing.T) (*isPkg.I, func()) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	_transactionModel.accounts = accountMock
	return is, func() {
		_transactionModel.accounts = nil
	}
}

func timeStampMock(month int) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(time.Date(2019, time.Month(month), 17, 16, 15, 14, 0, time.UTC))
	return ts
}
