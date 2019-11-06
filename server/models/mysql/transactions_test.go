package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"
	"testing"
)

func TestTransactionModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)

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
				Account: &models.Account{
					ID: 1,
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
			db, teardown := initializeForTransactions(t)
			defer teardown()

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

			got := models.Transaction{Account: &models.Account{}}

			stmt := `SELECT id, new_saldo, old_saldo, amount,created, account_id from transactions WHERE id=?`
			err = db.QueryRow(stmt, 1).Scan(&got.ID, &got.NewSaldo, &got.OldSaldo, &got.Amount, &got.Created, &got.Account.ID)
			is.NoErr(err)

			is.Equal(got.ID, tt.want.ID)                 // id does not match
			is.Equal(got.OldSaldo, tt.want.OldSaldo)     // oldSaldo does not match
			is.Equal(got.NewSaldo, tt.want.NewSaldo)     // newSaldo does not match
			is.True(!got.Created.IsZero())               // created is zero, should be timestamp
			is.Equal(got.Account.ID, tt.want.Account.ID) // accountId does not match

		})
	}
}

func initializeForTransactions(t *testing.T) (*sql.DB, func()) {
	return dbInitialized(t, "../testdata/transaction.sql")
}
