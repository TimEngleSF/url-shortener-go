package models

import (
	"time"

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

func (m *UserModel) Insert(name, email, password string) (*User, error) {
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
