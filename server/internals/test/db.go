package test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/internals/database"
)

// testDb will hold a db connection for multiple usages
type testDb struct {
	db             *sql.DB
	teardownScript string
}

// initDb connects to the database. will skip the calling test if TEST_DB_DSN is not set
// or connection to database is unsuccessful
func (tdb testDb) initDb(t *testing.T) *sql.DB {
	if tdb.db != nil {
		return tdb.db
	}

	dsn := "${TEST_DB_NAME}:${TEST_DB_NAME}@tcp(${TEST_DB_NAME}-test.db)/${TEST_DB_NAME}?parseTime=true&multiStatements=true"
	if os.Getenv("TEST_DB_DSN") != "" {
		dsn = os.Getenv("TEST_DB_DSN")
	}

	db, err := database.OpenDatabase(dsn)
	if err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}
	// db.SetConnMaxLifetime(1 * time.Second)
	// db.SetMaxOpenConns(100)
	return db
}

var testDBconn testDb

// getDb connects to database and returns db (*sql.DB) a setup function and a teardown function
// setup calls migration up scripts to bring database to latest and can
// receive filenames (strings) for additional setup scripts (sql-files)
// teardown calls migration down scripts
//
// setup and teardown can be called multiple times per test
func GetDb(t *testing.T, dbName, teardownScriptFile string, migrationDir string) (db *sql.DB, setup func(...string), teardown func()) {
	t.Helper()
	db = testDBconn.initDb(t)

	if err := db.Ping(); err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	testDBconn.teardownScript = func() string {
		content, err := ioutil.ReadFile(teardownScriptFile)
		if err != nil {
			t.Skipf("could not load teardown script: %v", err)
		}
		return string(content)
	}()

	teardown = func() {
		_, err := db.Exec(testDBconn.teardownScript)
		if err != nil {
			t.Fatalf("got error from teardown: %v", err)
		}
	}

	setup = func(setupScriptFileNames ...string) {
		if err := database.UpdateDatabase(db, dbName, migrationDir, true); err != nil {
			t.Skipf("could not migrate database, err: %v", err)
		}
		for _, filename := range setupScriptFileNames {
			setupScript, _ := ioutil.ReadFile(filename)
			_, err := db.Exec(string(setupScript))
			if err != nil {
				teardown()
				t.Fatalf("got error initializing script %q: %v", filename, err)
			}
		}
	}

	return db, setup, teardown
}
