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

func isIntegrationTest(t *testing.T) {
	if getEnvWithDefault("RUN_INTEGRATION", "0") == "0" {
		t.Skipf("skipping integration tests")
	}
}

func getEnvWithDefault(envName, defaultVal string) string {
	if val, ok := os.LookupEnv(envName); ok {
		return os.ExpandEnv(val)
	}
	return defaultVal
}

type testDb struct {
	db *sql.DB
}

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

func getTestDb(t *testing.T) (*sql.DB, func(...string), func()) {
	t.Helper()
	db := testDBconn.initDb(t)
	db.SetConnMaxLifetime(10 * time.Second)

	if err := db.Ping(); err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	migrationDir := getEnvWithDefault("DB_MIGRATIONS_DIR", defaultMigrationDir)
	dbName := getEnvWithDefault("TEST_DB_NAME", defaultDbName)

	teardownF := func() {
		err := database.DowngradeDatabase(db, dbName, migrationDir)
		if err != nil {
			t.Errorf("got error from teardown: %v", err)
		}
	}

	setupF := func(setupScriptFileNames ...string) {
		if err := database.UpdateDatabase(db, dbName, migrationDir, false); err != nil {
			t.Skipf("could not migrate database, err: %v", err)
		}
		for _, filename := range setupScriptFileNames {
			setupScript, _ := ioutil.ReadFile(filename)
			_, err := db.Exec(string(setupScript))
			if err != nil {
				teardownF()
				t.Fatalf("got error initializing script %q: %v", filename, err)
			}
		}
	}

	return db, setupF, teardownF
}
