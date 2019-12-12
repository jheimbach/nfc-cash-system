package mysql

import (
	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	isPkg "github.com/matryer/is"
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
				Amount:    -6,
				OldSaldo:  12,
				NewSaldo:  6,
				AccountId: 1,
			},
			want: &api.Transaction{
				Id:       1,
				OldSaldo: 12,
				NewSaldo: 6,
				Amount:   -6,
				Account: &api.Account{
					Id:        1,
					Name:      "testaccount",
					Saldo:     12,
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
				Amount:    -6,
				OldSaldo:  12,
				NewSaldo:  6,
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

			got, err := model.Create(tt.input.Amount, tt.input.OldSaldo, tt.input.NewSaldo, tt.input.AccountId)

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

func TestTransactionModel_GetAll(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("get all transactions", func(t *testing.T) {
		is := is.New(t)
		dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
		defer dbTeardown()

		model := TransactionModel{
			db:       db,
			accounts: NewAccountModel(db, NewGroupModel(db)),
		}
		transactions, err := model.GetAll()
		is.NoErr(err)

		is.Equal(len(transactions), 9) // expected 9 transactions

	})

	t.Run("no transactions found", func(t *testing.T) {
		dbSetup()
		defer dbTeardown()

		model := TransactionModel{
			db:       db,
			accounts: NewAccountModel(db, NewGroupModel(db)),
		}
		transactions, err := model.GetAll()
		is.NoErr(err)

		is.Equal(len(transactions), 0) // expected 0 transactions
	})

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

func TestTransactionModel_GetAllByAccount(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("get all transactions for account 1", func(t *testing.T) {
		is := is.New(t)

		dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
		defer dbTeardown()

		model := TransactionModel{
			db:       db,
			accounts: NewAccountModel(db, NewGroupModel(db)),
		}
		transactions, err := model.GetAllByAccount(1)
		is.NoErr(err)
		is.Equal(len(transactions), 5) // expected 5 transactions

		account := &api.Account{
			Id:        1,
			Name:      "testaccount",
			Saldo:     12,
			NfcChipId: "testchipid",
			Group: &api.Group{
				Id:   1,
				Name: "testgroup1",
			},
		}

		wantList := []*api.Transaction{
			{
				Id:       5,
				OldSaldo: 100,
				NewSaldo: 105,
				Amount:   -5,
				Created:  timeStampMock(5),
				Account:  account,
			},
			{
				Id:       4,
				OldSaldo: 105,
				NewSaldo: 100,
				Amount:   5,
				Created:  timeStampMock(4),
				Account:  account,
			},
			{
				Id:       3,
				OldSaldo: 110,
				NewSaldo: 105,
				Amount:   5,
				Created:  timeStampMock(3),
				Account:  account,
			},
			{
				Id:       2,
				OldSaldo: 115,
				NewSaldo: 110,
				Amount:   5,
				Created:  timeStampMock(2),
				Account:  account,
			},
			{
				Id:       1,
				OldSaldo: 120,
				NewSaldo: 115,
				Amount:   5,
				Created:  timeStampMock(1),
				Account:  account,
			},
		}

		is.Equal(transactions, wantList)
	})

	t.Run("no transactions found for account id 100", func(t *testing.T) {
		dbSetup()
		defer dbTeardown()

		model := TransactionModel{
			db: db,
		}

		transactions, err := model.GetAllByAccount(100)
		is.NoErr(err)

		is.Equal(len(transactions), 0) // expected 0 transactions
	})
}

func timeStampMock(month int) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(time.Date(2019, time.Month(month), 17, 16, 15, 14, 0, time.UTC))
	return ts
}
