package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/internals/database"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	defaultDbName       = "nfc-cash-system_test"
	defaultMigrationDir = "../migrations"
)

// isIntegrationTest will check if envvar RUN_INTEGRATION is non zero.
// if RUN_INTEGRATION is zero (or not set), tests that call isIntegrationTest will be skipped
func isIntegrationTest(t *testing.T) {
	if getEnvWithDefault("RUN_INTEGRATION", "0") == "0" {
		t.Skipf("skipping integration tests")
	}
}

// genEnvWithDefault is a helper function to get envvar or default value
func getEnvWithDefault(envName, defaultVal string) string {
	if val, ok := os.LookupEnv(envName); ok {
		return os.ExpandEnv(val)
	}
	return defaultVal
}

// testDb will hold a db connection for multiple usages
type testDb struct {
	db *sql.DB
}

// initDb connects to the database. will skip the calling test if TEST_DB_DSN is not set
// or connection to database is unsuccessful
func (tdb testDb) initDb(t *testing.T) *sql.DB {
	if tdb.db != nil {
		return tdb.db
	}

	dsn := getEnvWithDefault("TEST_DB_DSN", "")
	if dsn == "" {
		t.Skipf("no database dsn found, skipping test")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}
	return db
}

var testDBconn testDb

// getTestDb connects to database and returns db (*sql.DB) a setup function and a teardown function
// setup calls migration up scripts to bring database to latest and can
// receive filenames (strings) for additional setup scripts (sql-files)
// teardown calls migration down scripts
//
// setup and teardown can be called multiple times per test
func getTestDb(t *testing.T) (db *sql.DB, setup func(...string), teardown func()) {
	t.Helper()
	db = testDBconn.initDb(t)
	db.SetConnMaxLifetime(10 * time.Second)

	if err := db.Ping(); err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	migrationDir := getEnvWithDefault("DB_MIGRATIONS_DIR", defaultMigrationDir)
	dbName := getEnvWithDefault("TEST_DB_NAME", defaultDbName)

	teardown = func() {
		err := database.DowngradeDatabase(db, dbName, migrationDir)
		if err != nil {
			t.Errorf("got error from teardown: %v", err)
		}
	}

	setup = func(setupScriptFileNames ...string) {
		if err := database.UpdateDatabase(db, dbName, migrationDir, false); err != nil {
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
