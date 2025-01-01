package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO users (name, email, hashed_password, created)
  VALUES($1, $2, $3, NOW() AT TIME ZONE 'UTC');`

	_, err = model.DB.Exec(sqlStatement, name, email, string(hashedPassword))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (model *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (model *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
