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
	var id int
	var hashedPassword []byte
	sqlStatement := "SELECT id, hashed_password FROM users WHERE email = $1"
	err := model.DB.QueryRow(sqlStatement, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (model *UserModel) Exists(id int) (bool, error) {
	var exists bool
	sqlStatement := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"
	err := model.DB.QueryRow(sqlStatement, id).Scan(&exists)
	return exists, err
}
