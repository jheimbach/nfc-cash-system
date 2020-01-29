package mysql

import (
	"database/sql"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test"
)

var (
	_userModel        *UserRepository
	_groupModel       *GroupRepository
	_accountModel     *AccountRepository
	_transactionModel *TransactionRepository
	_conn             *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	_conn, _, err = test.DbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer _conn.Close()

	_userModel = NewUserModel(_conn)
	_groupModel = NewGroupRepository(_conn)
	_accountModel = NewAccountRepository(_conn, nil)
	_transactionModel = NewTransactionRepository(_conn, nil)

	os.Exit(m.Run())
}

func dataFor(t string) string {
	return path.Join("./testdata", strings.Join([]string{t, "sql"}, "."))
}

func teardownDB(db *sql.DB) func() error {
	return test.TeardownDB(db, dataFor("teardown"))
}
