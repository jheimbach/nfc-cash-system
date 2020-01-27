package mysql

import (
	"database/sql"
	"path"
	"strings"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
)

const (
	DbName       = "nfc-cash-system_test"
	MigrationDir = "../../migrations"
)

func dataFor(t string) string {
	return path.Join("../../testdata", strings.Join([]string{t, "sql"}, "."))
}

func getTestDb(t *testing.T) (*sql.DB, func(...string), func()) {
	migrationDir := test.EnvWithDefault("DB_MIGRATIONS_DIR", MigrationDir)
	dbName := test.EnvWithDefault("TEST_DB_NAME", DbName)
	return test.GetDb(t, dbName, dataFor("teardown"), migrationDir)
}
