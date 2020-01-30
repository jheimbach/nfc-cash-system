package test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/docker/go-connections/nat"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/database"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	MysqlUser           = "test"
	MysqlPassword       = "test"
	MysqlDatabase       = "mysql_test"
	defaultMigrationDir = "./migrations"
)

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

func DbConnection(migrationsPath string) (db *sql.DB, teardown func(), err error) {
	addr, _, err := StartDbContainer("")
	db, err = OpenAndMigrateDatabase(addr, migrationsPath)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}

func OpenAndMigrateDatabase(addr string, migrationDir string) (*sql.DB, error) {
	db, err := database.OpenDatabase(database.CreateDsn(MysqlUser, MysqlPassword, addr, MysqlDatabase))
	if err != nil {
		return nil, fmt.Errorf("could not conntect to database: %w", err)
	}

	if migrationDir == "" {
		migrationDir = defaultMigrationDir
	}

	err = database.UpdateDatabase(db, MysqlDatabase, migrationDir, false)
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
