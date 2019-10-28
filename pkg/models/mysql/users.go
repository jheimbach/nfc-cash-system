package mysql

import (
	"database/sql"
	"github.com/JHeimbach/nfc-cash-system/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	db *sql.DB
}

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

func (u UserModel) Get(id int) (*models.User, error) {
	m := &models.User{}

	getSql := `SELECT id,name,email,created FROM users WHERE id = ?`
	err := u.db.QueryRow(getSql, id).Scan(&m.ID, &m.Name, &m.Email, &m.Created)
	return m, err
}
