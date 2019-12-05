package mysql

import (
	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
	"testing"
)

func TestTransactionModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		amount    float64
		oldSaldo  float64
		newSaldo  float64
		accountId int
	}
	tests := []struct {
		name        string
		args        args
		want        *models.Transaction
		wantErr     bool
		expectedErr error
	}{
		{
			name: "create new transaction",
			args: args{
				amount:    -6,
				oldSaldo:  12,
				newSaldo:  6,
				accountId: 1,
			},
			want: &models.Transaction{
				ID:       1,
				OldSaldo: 12,
				NewSaldo: 6,
				Amount:   -6,
				Account: api.Account{
					Id: 1,
				},
			},
		}, {
			name: "create new transaction with nonexistent account",
			args: args{
				amount:    -6,
				oldSaldo:  12,
				newSaldo:  6,
				accountId: 100,
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
				db: db,
			}

			err := model.Create(tt.args.amount, tt.args.oldSaldo, tt.args.newSaldo, tt.args.accountId)

			if tt.wantErr {
				if err != tt.expectedErr {
					t.Errorf("got err %v, expected %v", err, tt.expectedErr)
				}
				return
			}

			is.NoErr(err)

			got := models.Transaction{Account: api.Account{}}

			stmt := `SELECT id, new_saldo, old_saldo, amount,created, account_id from transactions WHERE id=?`
			err = db.QueryRow(stmt, 1).Scan(
				&got.ID, &got.NewSaldo, &got.OldSaldo, &got.Amount, &got.Created, &got.Account.Id,
			)
			is.NoErr(err)

			is.Equal(got.ID, tt.want.ID)                 // id does not match
			is.Equal(got.OldSaldo, tt.want.OldSaldo)     // oldSaldo does not match
			is.Equal(got.NewSaldo, tt.want.NewSaldo)     // newSaldo does not match
			is.True(!got.Created.IsZero())               // created is zero, should be timestamp
			is.Equal(got.Account.Id, tt.want.Account.Id) // accountId does not match

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
			db: db,
		}
		transactions, err := model.GetAll()
		is.NoErr(err)

		is.Equal(len(transactions), 9) // expected 9 transactions

	})

	t.Run("no transactions found", func(t *testing.T) {
		dbSetup()
		defer dbTeardown()

		model := TransactionModel{
			db: db,
		}
		transactions, err := model.GetAll()
		is.NoErr(err)

		is.Equal(len(transactions), 0) // expected 0 transactions
	})

}

func TestTransactionModel_GetAllPaged(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := TransactionModel{
		db: db,
	}

	page1, err := model.GetAllPaged(1, 5)
	is.NoErr(err)

	is.Equal(page1.CurrentPage, 1)       // currentpage should be 1
	is.Equal(page1.MaxPage, 2)           // maxpage should be 2
	is.Equal(len(page1.Transactions), 5) // expected 5 transactions

	page2, err := model.GetAllPaged(2, 5)
	is.NoErr(err)

	is.Equal(page2.CurrentPage, 2)       // currentpage should be 2
	is.Equal(page2.MaxPage, 2)           // maxpage should be 2
	is.Equal(len(page2.Transactions), 4) // expected 4 transactions
}

func TestTransactionModel_GetAllByAccount(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("get all transactions for accountId 1", func(t *testing.T) {
		is := is.New(t)

		dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
		defer dbTeardown()

		model := TransactionModel{
			db: db,
		}
		transactions, err := model.GetAllByAccount(1)
		is.NoErr(err)

		is.Equal(len(transactions), 5) // expected 5 transactions
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

func TestTransactionModel_GetAllByAccountPaged(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

	db, dbSetup, dbTeardown := getTestDb(t)
	dbSetup("../testdata/transaction.sql", "../testdata/transaction_list.sql")
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := TransactionModel{
		db: db,
	}

	page1, err := model.GetAllByAccountPaged(1, 1, 3)
	is.NoErr(err)

	is.Equal(page1.CurrentPage, 1)       // currentpage should be 1
	is.Equal(page1.MaxPage, 2)           // maxpage should be 2
	is.Equal(len(page1.Transactions), 3) // expected 5 transactions

	page2, err := model.GetAllByAccountPaged(1, 2, 3)
	is.NoErr(err)

	is.Equal(page2.CurrentPage, 2)       // currentpage should be 2
	is.Equal(page2.MaxPage, 2)           // maxpage should be 2
	is.Equal(len(page2.Transactions), 2) // expected 2 transactions

}
