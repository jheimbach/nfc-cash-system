package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserModel handles the users from the database
type UserModel struct {
	db *sql.DB
}

// Insert creates a new user in the database.
// if a user with the same email already exists, Insert will return a models.ErrDuplicateEmail
func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	insertSql := `INSERT INTO users (name,email,hashed_password,created) VALUES(?,?,?,UTC_TIMESTAMP())`
	_, err = u.db.Exec(insertSql, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return models.ErrDuplicateEmail
			}
		}
	}

	return err
}

// Get returns a models.User from given id, if id does not exists Get will return a models.ErrNotFound
func (u *UserModel) Get(id int) (*models.User, error) {
	m := &models.User{}

	getSql := `SELECT id,name,email,created FROM users WHERE id = ?`
	err := u.db.QueryRow(getSql, id).Scan(&m.ID, &m.Name, &m.Email, &m.Created)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return m, nil
}

// Authenticate returns the id for a user if it exists with given email and password
// if email does not exists or the password is wrong Authenticate will return a models.ErrInvalidCredentials
func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := u.db.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}
