package qr

import (
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeInterface interface {
	CreateMedium(url string) (png []byte, err error)
}

type QRCode struct{}

func (qr *QRCode) CreateMedium(url string) (png []byte, err error) {
	timestamp := time.Now().Format(time.RFC850)
	err = qrcode.WriteFile(url, qrcode.Medium, 256, timestamp+".png")
	if err != nil {
		return nil, err
	}

	return
}
