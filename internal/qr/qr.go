package qr

import (
	"strings"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeInterface interface {
	CreateMedium(url string) (path string, err error)
}

type QRCode struct{}

func (qr *QRCode) CreateMedium(url string) (path string, err error) {
	id := uuid.New()
	path = "./ui/static/qr/" + id.String() + ".png"
	err = qrcode.WriteFile(url, qrcode.Medium, 256, path)
	if err != nil {
		return "", err
	}

	path, _ = strings.CutPrefix(path, "./ui")
	return path, nil
}
