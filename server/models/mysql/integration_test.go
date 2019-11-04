package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/internals/database"
	"os"
	"testing"
)

const (
	testDatabaseName = "nfc-cash-system_test"
	migrationDbFiles = "../migrations"
)

func getTestDb(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
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

	if err = database.UpdateDatabase(db, testDatabaseName, migrationDbFiles, false); err != nil {
		t.Skipf("could not migrate database, err: %v", err)
	}

	teardownF := func() {
		err := database.DowngradeDatabase(db, testDatabaseName, migrationDbFiles)
		if err != nil {
			t.Log(err)
		}
		db.Close()
	}

	return db, teardownF
}
