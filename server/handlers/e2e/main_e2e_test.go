package e2e

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
)

var (
	_conn         *sql.DB = nil
	_restEndpoint         = ":"
	_aTkn                 = ""
	_rTkn                 = ""
	_dbinitialized = false
)

func TestMain(m *testing.M) {
	var err error
	var teardown func()
	_conn, _restEndpoint, teardown, err = startServers()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		teardown()
		_conn.Close()
	}()

	err = test.SetupDB(_conn, dataFor("end-to-end"))
	if err != nil {
		log.Fatal(err)
	}
	_dbinitialized = true
	defer test.TeardownDB(_conn, dataFor("teardown"))()

	_aTkn, _rTkn, err = login()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
