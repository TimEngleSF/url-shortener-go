package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

var mockLink = models.Link{
	ID:          1,
	RedirectUrl: "https://google.com",
	Suffix:      "abc123",
	CreatedAt:   time.Now(),
	ShortUrl:    "",
}

type LinkMock struct{}

func (m *LinkMock) Insert(ctx context.Context, redirectUrl, suffix string) (models.Link, error) {
	l := models.Link{
		ID:          2,
		RedirectUrl: redirectUrl,
		Suffix:      suffix,
		CreatedAt:   time.Now(),
	}

	shortUrl, err := l.CreateShortUrl("example.com")
	if err != nil {
		return models.Link{}, err
	}

	l.ShortUrl = shortUrl

	return l, nil
}

func (m *LinkMock) GetBySuffix(ctx context.Context, suffix string) (models.Link, error) {
	fmt.Println(suffix)
	fmt.Println(suffix == "abc123")
	switch suffix {
	case "abc123":
		return mockLink, nil
	default:
		return models.Link{}, models.ErrNoRecord
	}
}

func (m *LinkMock) GetByURL(ctx context.Context, url string) (models.Link, error) {
	return models.Link{}, nil
}

func (m *LinkMock) URLExists(urlStr string) (bool, error) {
	return false, nil
}
