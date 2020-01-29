package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const DefaultDSN = "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true&multiStatements=true"

func OpenDatabase(dsn string) (*sql.DB, error) {
	populated := os.ExpandEnv(dsn)
	db, err := sql.Open("mysql", populated)
	if err != nil {
		return nil, err
	}

	if err := ping(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		time.Sleep(1 * time.Second)
		return repeatedPing(db, 1)
	}
	return nil
}

func repeatedPing(db *sql.DB, times int) error {
	if err := db.Ping(); err != nil {
		if times == 10 {
			return err
		}
		time.Sleep(1 * time.Second)
		return repeatedPing(db, times+1)
	}
	return nil
}

func CheckEnvVars() error {
	vars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_NAME"}
	for _, s := range vars {
		_, ok := os.LookupEnv(s)
		if !ok {
			return fmt.Errorf("envvar %s missing, but it is required", s)
		}
	}

	return nil
}
