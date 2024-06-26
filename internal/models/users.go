package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	Insert(ctx context.Context, name, email, password string) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
	Get(ctx context.Context, id int) (*User, error)
	ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error
	ExistsByID(ctx context.Context, id int) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	AddLink(ctx context.Context, user_id, link_id int) error
	HasLink(ctx context.Context, user_id, link_id int) (bool, error)
	GetLinks(ctx context.Context, user_id int, host string) ([]Link, error)
}

func (m *UserModel) Insert(ctx context.Context, name, email, password string) error {
	stmt := `
  INSERT INTO users (
    name, 
    email, 
    password,
    created_at, 
    updated_at
  ) 
  VALUES (
    $1, 
    $2,
    $3,
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC', 
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
  )
  RETURNING user_id;
  `

	var id int
	lowEmail := strings.ToLower(email)
	hashedPass, err := HashPassword(password)
	if err != nil {
		return err
	}

	err = m.DB.QueryRow(ctx, stmt, name, lowEmail, hashedPass).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		}

		return err
	}
	return nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (*User, error) {
	stmt := `
  SELECT user_id, name, email, password
  FROM users
  WHERE $1 = email
  `
	var user User
	err := m.DB.QueryRow(ctx, stmt, strings.ToLower(email)).
		Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, ErrInvalidCredentials
			}
		}
		return nil, err
	}

	validPass := CheckPasswordHash(password, string(user.HashedPassword))
	if !validPass {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

func (m *UserModel) Get(ctx context.Context, id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	return nil
}

func (m *UserModel) ExistsByID(ctx context.Context, id int) (bool, error) {
	doesExist, err := exists(m.DB, ctx, id)
	if err != nil {
		return false, err
	}

	return doesExist, nil
}

func (m *UserModel) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	lowEmail := strings.ToLower(email)
	doessExist, err := exists(m.DB, ctx, lowEmail)
	if err != nil {
		return false, err
	}
	return doessExist, nil
}

func (m *UserModel) AddLink(ctx context.Context, user_id, link_id int) error {
	stmt := `
  INSERT INTO user_links (
    user_id,
    link_id
  )
  VALUES (
    $1,
    $2
  )
  `
	_, err := m.DB.Exec(ctx, stmt, user_id, link_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetLinks(ctx context.Context, user_id int, host string) ([]Link, error) {
	stmt := `
  SELECT l.*
  FROM links l
  JOIN users_links ul ON l.link_id = ul.link_id
  WHERE ul.user_id = $1
  `
	links := []Link{}

	rows, err := m.DB.Query(ctx, stmt, user_id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Code)
			fmt.Println(pgErr.Message)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var l Link
		err := rows.Scan(&l.ID, &l.RedirectUrl, &l.Suffix, &l.QRUrl, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		l.ShortUrl, _ = l.CreateShortUrl(host)
		links = append(links, l)
	}
	return links, nil
}

func (m *UserModel) HasLink(ctx context.Context, user_id, link_id int) (bool, error) {
	stmt := `
  SELECT EXISTS(
    SELECT 1 FROM user_links
    WHERE user_id = $1 AND link_id = $2);
  `
	var exists bool
	err := m.DB.QueryRow(ctx, stmt, user_id, link_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Check if there is a row based on an email or id value
func exists(db *pgxpool.Pool, ctx context.Context, value interface{}) (bool, error) {
	var targetColumn string
	switch value.(type) {
	case int:
		targetColumn = "user_id"
	case string:
		targetColumn = "email"
	default:
		return false, errors.New("value argument invalid type")
	}

	stmt := fmt.Sprintf(
		`SELECT 
    EXISTS(
      SELECT 1
      FROM users
      WHERE %s = $1 
    )
`, targetColumn)

	var exists bool
	err := db.QueryRow(ctx, stmt, value).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) InsertLink(user_id, link_id int) error {
	return nil
}
