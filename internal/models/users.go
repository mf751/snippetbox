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

// Becasue the application struct is to expect this
// so that it can recive both the usersModel and the mocks.usersModel
type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	GetAccountInfo(id int) (User, error)
	ChangePassword(id int, currentPassword, newPassword string) error
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

func (model *UserModel) GetAccountInfo(id int) (User, error) {
	user := User{}
	sqlStatement := `SELECT name, email, created FROM users WHERE id=$1`
	err := model.DB.QueryRow(sqlStatement, id).Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (model *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte
	sqlStatement1 := `SELECT hashed_password FROM users WHERE id=$1;`
	err := model.DB.QueryRow(sqlStatement1, id).Scan(&hashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	sqlStatement2 := `UPDATE users SET hashed_password=$1 WHERE id=$2`
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	if _, err = model.DB.Exec(sqlStatement2, string(newHashedPassword), id); err != nil {
		return err
	}

	return nil
}
