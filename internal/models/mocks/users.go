package mocks

import (
	"context"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

type UserMock struct{}

func (m *UserMock) Insert(ctx context.Context, name, email, password string) error {
	return nil
}

func (m *UserMock) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	return nil, nil
}

func (m *UserMock) Get(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}

func (m *UserMock) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	return nil
}

func (m *UserMock) ExistsByID(ctx context.Context, id int) (bool, error) {
	return false, nil
}

func (m *UserMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}
