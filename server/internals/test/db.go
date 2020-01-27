package test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/internals/database"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	MysqlUser     = "test"
	MysqlPassword = "test"
	MysqlDatabase = "mysql_test"
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

func createDSN(server, port string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", MysqlUser, MysqlPassword, server, port, MysqlDatabase)
}

func StartDbContainer(networkName string) (local, network string, err error) {
	ctx := context.Background()
	mysqlPort, err := nat.NewPort("tcp", "3306")

	if err != nil {
		return "", "", err

	}

	containerReq := tc.ContainerRequest{
		Image:        "mysql",
		ExposedPorts: []string{mysqlPort.Port()},
		Env: map[string]string{
			"MYSQL_RANDOM_ROOT_PASSWORD": "1",
			"MYSQL_USER":                 MysqlUser,
			"MYSQL_PASSWORD":             MysqlPassword,
			"MYSQL_DATABASE":             MysqlDatabase,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
			wait.ForListeningPort(mysqlPort),
		),
	}
	if networkName != "" {
		containerReq.Networks = []string{networkName}
		containerReq.NetworkAliases = map[string][]string{
			networkName: {"mysql-server"},
		}
	}

	mysqlC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})

	if err != nil {

		return "", "", fmt.Errorf("could not start container: %w", err)
	}



	mappedPort, err := mysqlC.MappedPort(ctx, mysqlPort)
	if err != nil {
		return "", "", fmt.Errorf("could not get internal endpoint address: %w", err)
	}

	return net.JoinHostPort("localhost", mappedPort.Port()), net.JoinHostPort("mysql-server", mappedPort.Port()), nil
}

func DbConnection() (db *sql.DB, teardown func(), err error) {
	addr, _, err := StartDbContainer("")
	db, err = OpenAndMigrateDatabase(addr)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}

func OpenAndMigrateDatabase(addr string) (*sql.DB, error) {
	a := strings.Split(addr, ":")
	db, err := database.OpenDatabase(createDSN(a[0], a[1]))
	if err != nil {
		return nil, fmt.Errorf("could not conntect to database: %w", err)
	}

	err = database.UpdateDatabase(db, MysqlDatabase, "../../migrations", false)
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}
	return db, nil
}

func SetupDB(db *sql.DB, setupScriptFileNames ...string) error {
	for _, filename := range setupScriptFileNames {
		setupScript, _ := ioutil.ReadFile(filename)
		_, err := db.Exec(string(setupScript))
		if err != nil {
			return fmt.Errorf("initializing script %q: %v", filename, err)
		}
	}
	return nil
}

func TeardownDB(db *sql.DB, teardownScriptFileName string) func() error {
	return func() error {
		teardownScript, _ := ioutil.ReadFile(teardownScriptFileName)
		_, err := db.Exec(string(teardownScript))
		if err != nil {
			return fmt.Errorf("got error from teardown: %w", err)
		}
		return nil
	}
}

