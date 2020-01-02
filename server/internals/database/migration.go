package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initMigrationDriver(db *sql.DB, databaseName, migrationDir string) (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationDir), databaseName, driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func UpdateDatabase(db *sql.DB, databaseName, migrationDir string, force bool) error {
	m, err := initMigrationDriver(db, databaseName, migrationDir)
	if err != nil {
		return err
	}

	if force {
		err = forcing(m)
		if err != nil {
			return err
		}
	}

	err = m.Up()

	return checkMigrationError(err)
}

func DowngradeDatabase(db *sql.DB, databaseName, migrationDir string) error {
	m, err := initMigrationDriver(db, databaseName, migrationDir)
	if err != nil {
		return err
	}

	err = m.Down()

	return checkMigrationError(err)
}

func forcing(m *migrate.Migrate) error {
	version, isDirty, err := m.Version()
	if err != nil {
		return err
	}

	if !isDirty {
		return nil
	}

	err = m.Force(int(version))
	if err != nil {
		return err
	}

	return nil
}

func checkMigrationError(err error) error {
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
