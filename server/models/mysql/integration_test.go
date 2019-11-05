package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/internals/database"
	"os"
	"testing"
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

func getTestDb(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	dsn := getEnvWithDefault("TEST_DB_DSN", "")
	if dsn == "" {
		t.Skipf("no database dsn found, skipping test")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	if err = db.Ping(); err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	migrationDir := getEnvWithDefault("DB_MIGRATIONS_DIR", defaultMigrationDir)
	dbName := getEnvWithDefault("TEST_DB_NAME", defaultDbName)

	if err = database.UpdateDatabase(db, dbName, migrationDir, false); err != nil {
		t.Skipf("could not migrate database, err: %v", err)
	}

	teardownF := func() {
		err := database.DowngradeDatabase(db, dbName, migrationDir)
		if err != nil {
			t.Log(err)
		}
		db.Close()
	}

	return db, teardownF
}
