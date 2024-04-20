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

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
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

func (m *UserModel) ChangeName(id int, name string) error {
	return nil
}
