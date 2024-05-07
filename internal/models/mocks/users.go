package mocks

import (
	"context"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

type UserMock struct{}

func (m *UserMock) Insert(ctx context.Context, name, email, password string) error {
	switch email {
	case "dupe@email.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserMock) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	if email == "john@email.com" && password == "pa$$word" {
		return &models.User{ID: 1}, nil
	}
	return &models.User{ID: 0}, models.ErrInvalidCredentials
}

func (m *UserMock) Get(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}

func (m *UserMock) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	return nil
}

func (m *UserMock) ExistsByID(ctx context.Context, id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	switch email {
	case "john@email.com":
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserMock) AddLink(ctx context.Context, user_id, link_id int) error {
	return nil
}

func (m *UserMock) HasLink(ctx context.Context, user_id, link_id int) (bool, error) {
	return true, nil
}

func (m *UserMock) GetLinks(ctx context.Context, user_id int, host string) ([]models.Link, error) {
	return nil, nil
}
