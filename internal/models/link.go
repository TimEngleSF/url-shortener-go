package models

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var SUFFIX_LENGTH = 6

type LinkModel struct {
	DB *pgxpool.Pool
}

type Link struct {
	ID          *int
	RedirectUrl string
	Suffix      string
	CreatedAt   *time.Time
	ShortUrl    string
}

func (m *LinkModel) Insert(ctx context.Context, redirectUrl, suffix string) (Link, error) {
	stmt := `INSERT INTO links (redirect_url, suffix) VALUES ($1, $2)`
	_, err := m.DB.Exec(ctx, stmt, redirectUrl, suffix)
	if err != nil {
		return Link{}, err
	}
	return Link{RedirectUrl: redirectUrl, Suffix: suffix}, nil
}

func (m *LinkModel) View(ctx context.Context, suffix string) (Link, error) {
	var link Link
	stmt := `SELECT id, redirect_url, suffix, created_at FROM links
  WHERE suffix = $1`
	err := m.DB.QueryRow(ctx, stmt, suffix).Scan(&link.ID, &link.RedirectUrl, &link.Suffix, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Link{}, ErrNoRecord
		} else {
			return Link{}, err
		}
	}
	return link, nil
}

func (m *LinkModel) SuffixExists(suffix string) (bool, error) {
	return false, nil
}

func (m *LinkModel) RedirectUrlExists(urlStr string) (bool, error) {
	return false, nil
}

func CreateSuffix() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	suffix := make([]byte, SUFFIX_LENGTH)
	for i := range suffix {
		suffix[i] = charset[rand.Intn(len(charset))]
	}
	return string(suffix)
}

func (m *LinkModel) GetBySuffix(ctx context.Context, suffix string) (Link, error) {
	var link Link
	stmt := `SELECT id, redirect_url, suffix, created_at FROM links
  WHERE suffix = $1`
	err := m.DB.QueryRow(ctx, stmt, suffix).Scan(&link.ID, &link.RedirectUrl, &link.Suffix, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Link{}, ErrNoRecord
		} else {
			return Link{}, err
		}
	}
	return link, nil
}

func (m *LinkModel) GetByURL(ctx context.Context, url string) (Link, error) {
	var link Link
	stmt := `SELECT id, redirect_url, suffix, created_at FROM links
  WHERE redirect_url = $1`
	err := m.DB.QueryRow(ctx, stmt, url).Scan(&link.ID, &link.RedirectUrl, &link.Suffix, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Link{}, ErrNoRecord
		} else {
			return Link{}, err
		}
	}
	return link, nil
}

func (m *LinkModel) URLExists(urlStr string) (bool, error) {
	stmt := `SELECT COUNT(*) FROM links WHERE redirect_url = $1`
	var exists bool
	err := m.DB.QueryRow(context.Background(), stmt, urlStr).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (l Link) CreateShortUrl(host string) (string, error) {
	if l.Suffix == "" {
		return "", ErrEmptySuffix
	}
	return "https://" + host + "/" + l.Suffix, nil
}
