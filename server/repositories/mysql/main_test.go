package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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

var (
	_userModel        *UserRepository
	_groupModel       *GroupRepository
	_accountModel     *AccountRepository
	_transactionModel *TransactionRepository
	_conn             *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	mysqlPort, err := nat.NewPort("tcp", "3306")
	if err != nil {
		log.Fatal(err)
	}
	mysqlC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
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
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("start: %v", err)
	}

	ip, err := mysqlC.Host(ctx)
	if err != nil {
		log.Fatalf("getting host: %v", err)
	}
	port, err := mysqlC.MappedPort(ctx, mysqlPort)
	if err != nil {
		log.Fatalf("getting port: %v", err)
	}

	_conn, err = database.OpenDatabase(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", MysqlUser, MysqlPassword, ip, port.Port(), MysqlDatabase))
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	err = database.UpdateDatabase(_conn, MysqlDatabase, "../../migrations", false)
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	_userModel = NewUserModel(_conn)
	_groupModel = NewGroupRepository(_conn)
	_accountModel = NewAccountRepository(_conn, nil)
	_transactionModel = NewTransactionRepository(_conn, nil)

	os.Exit(m.Run())
}

func setupDB(db *sql.DB, setupScriptFileNames ...string) error {
	for _, filename := range setupScriptFileNames {
		setupScript, _ := ioutil.ReadFile(filename)
		_, err := db.Exec(string(setupScript))
		if err != nil {
			return fmt.Errorf("initializing script %q: %v", filename, err)
		}
	}
	return nil
}

func teardownDB(db *sql.DB) func() error {
	return func() error {
		teardownScript, _ := ioutil.ReadFile("../../testdata/teardown.sql")
		_, err := db.Exec(string(teardownScript))
		if err != nil {
			return fmt.Errorf("got error from teardown: %w", err)
		}
		return nil
	}
}

func dataFor(t string) string {
	return path.Join("../../testdata", strings.Join([]string{t, "sql"}, "."))
}
