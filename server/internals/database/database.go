package database

import (
	"database/sql"
	"fmt"
	"os"
)

const DefaultDSN = "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true&multiStatements=true"

func OpenDatabase(dsn string) (*sql.DB, error) {
	populated := os.ExpandEnv(dsn)
	db, err := sql.Open("mysql", populated)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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
