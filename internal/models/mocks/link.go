package mocks

import (
	"context"
	"time"

	"github.com/TimEngleSF/url-shortener-go/internal/models"
)

var testInsertedSuffix = "123abc"

var mockLink = models.Link{
	ID:          1,
	RedirectUrl: "https://google.com",
	Suffix:      "abc123",
	CreatedAt:   time.Now(),
	QRUrl:       "someS3Url",
	ShortUrl:    "",
}

type LinkMock struct{}

func (m *LinkMock) Insert(ctx context.Context, redirectUrl, suffix, qrUrl, host string) (models.Link, error) {
	l := models.Link{
		ID:          2,
		RedirectUrl: redirectUrl,
		Suffix:      testInsertedSuffix,
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
	switch suffix {
	case "abc123":
		return mockLink, nil
	default:
		return models.Link{}, models.ErrNoRecord
	}
}

func (m *LinkMock) GetByURL(ctx context.Context, url, host string) (models.Link, error) {
	switch url {
	case "https://google.com":
		mockLink.ShortUrl, _ = mockLink.CreateShortUrl(host)
		return mockLink, nil
	default:
		return models.Link{}, models.ErrNoRecord
	}
}

func (m *LinkMock) URLExists(urlStr string) (bool, error) {
	return false, nil
}
