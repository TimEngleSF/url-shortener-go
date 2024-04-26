package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
	Get(id int) (*User, error)
	ChangePassword(id int, currentPassword, newPassword string) error
	Exists(email string) (bool, error)
}

func (m *UserModel) Insert(ctx context.Context, name, email, password string) (*User, error) {
	stmt := `
  INSERT INTO users (
    name, 
    email, 
    password,
    created_at, 
    updated_at)
  VALUES(?,?,?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
    returning user_id`

	var id int
	err := m.DB.QueryRow(ctx, stmt, name, email, password).Scan(&id)
	if err != nil {
		// TODO: Check to make sure that the customized email error is returned by psql below
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message)
			fmt.Println(pgErr.Code)
		}

		return nil, err
	}
	fmt.Println(id)
	return nil, nil
}

func (m *UserModel) Authenticate(email, password string) (*User, error) {
	return nil, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	return nil
}

func (m *UserModel) Exists(email string) (bool, error) {
	return false, nil
}
