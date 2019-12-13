package mysql

import (
	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	isPkg "github.com/matryer/is"
	"sort"
	"testing"
	"time"
)

func TestTransactionModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()
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
			expectedErr: models.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dbSetup("../testdata/transaction.sql")
			defer dbTeardown()

			model := TransactionModel{
				db:       db,
				accounts: NewAccountModel(db, NewGroupModel(db)), //todo create mock
			}

			got, err := model.Create(tt.input.Amount, tt.input.AccountId)

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
			err = db.QueryRow(stmt, 1).Scan(
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
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
	defer db.Close()
	defer dbTeardown()

	t.Run("get a transaction", func(t *testing.T) {
		is := is.New(t)

		model := TransactionModel{
			db: db,
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

		transaction, err := model.Read(1)
		is.NoErr(err)

		is.Equal(transaction, want)
	})

	t.Run("no transactions found", func(t *testing.T) {
		is := is.New(t)

		model := TransactionModel{
			db: db,
		}
		transaction, err := model.Read(100)
		is.Equal(transaction, nil)

		if err != models.ErrNotFound {
			t.Errorf("got err %v, expected %v", err, models.ErrNotFound)
		}
	})

}

func TestTransactionModel_GetAll(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		limit, offset int32
		order         string
	}
	tests := []struct {
		name    string
		dbSetup []string
		input   args
		want    []*api.Transaction
	}{
		{
			name:    "get all transactions",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input:   args{},
			want:    transisitonList(0),
		},
		{
			name:  "no transactions found",
			input: args{},
		},
		{
			name:    "get transactions with limit",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				limit: 5,
			},
			want: transisitonList(0)[:5],
		},
		{
			name:    "get transactions with limit and offset",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				limit:  3,
				offset: 2,
			},
			want: transisitonList(0)[2:5],
		},
		{
			name:    "get all transactions with order DESC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				order: "DESC",
			},
			want: transisitonList(0),
		},
		{
			name:    "get all transactions with order desc",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				order: "desc",
			},
			want: transisitonList(0),
		},
		{
			name:    "get all transactions default order is DESC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				order: "something invalid",
			},
			want: transisitonList(0),
		},
		{
			name:    "get all transactions with order ASC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				order: "ASC",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(0))
				sort.Sort(s)
				return s
			}(),
		},
		{
			name:    "get all transactions with order asc",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				order: "asc",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(0))
				sort.Sort(s)
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dbSetup(tt.dbSetup...)
			defer dbTeardown()

			model := TransactionModel{
				db:       db,
				accounts: NewAccountModel(db, NewGroupModel(db)),
			}
			got, err := model.GetAll(tt.input.order, tt.input.limit, tt.input.offset)
			is.NoErr(err)
			is.Equal(got, tt.want)
		})
	}

}

func TestTransactionModel_GetAllByAccount(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		accountId, limit, offset int32
		order                    string
	}
	tests := []struct {
		name    string
		dbSetup []string
		input   args
		want    []*api.Transaction
	}{
		{
			name:    "get all transactions by account with id 1",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
			},
			want: transisitonList(1),
		},
		{
			name: "no transactions found",
			input: args{
				accountId: 1,
			},
		},
		{
			name:    "get transactions with limit",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				limit:     3,
			},
			want: transisitonList(1)[:3],
		},
		{
			name:    "get transactions with limit and offset",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				limit:     3,
				offset:    2,
			},
			want: transisitonList(1)[2:5],
		},
		{
			name:    "get transactions with order DESC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				order:     "DESC",
			},
			want: transisitonList(1),
		},
		{
			name:    "get transactions with order desc",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				order:     "desc",
			},
			want: transisitonList(1),
		},
		{
			name:    "get transactions default order is DESC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				order:     "something invalid",
			},
			want: transisitonList(1),
		},
		{
			name:    "get transactions with order ASC",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				order:     "ASC",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(1))
				sort.Sort(s)
				return s
			}(),
		},
		{
			name:    "get transactions with order asc",
			dbSetup: []string{"../testdata/transaction.sql", "../testdata/transaction_list.sql"},
			input: args{
				accountId: 1,
				order:     "asc",
			},
			want: func() []*api.Transaction {
				s := SortTransactions(transisitonList(1))
				sort.Sort(s)
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dbSetup(tt.dbSetup...)
			defer dbTeardown()

			model := TransactionModel{
				db:       db,
				accounts: NewAccountModel(db, NewGroupModel(db)),
			}
			got, err := model.GetAllByAccount(tt.input.accountId, tt.input.order, tt.input.limit, tt.input.offset)
			is.NoErr(err)
			is.Equal(got, tt.want)
		})
	}

}

func timeStampMock(month int) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(time.Date(2019, time.Month(month), 17, 16, 15, 14, 0, time.UTC))
	return ts
}

func transisitonList(accountId int) []*api.Transaction {
	accountOne := &api.Account{
		Id:        1,
		Name:      "testaccount",
		Saldo:     12,
		NfcChipId: "testchipid",
		Group: &api.Group{
			Id:   1,
			Name: "testgroup1",
		},
	}
	accountTwo := &api.Account{
		Id:        2,
		Name:      "testaccount1",
		Saldo:     120,
		NfcChipId: "testchipid2",
		Group: &api.Group{
			Id:   1,
			Name: "testgroup1",
		},
	}

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
