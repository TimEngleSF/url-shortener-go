package qr

import "fmt"

type QRCodeMock struct {
	CreateMockCalled bool
}

func (mock *QRCodeMock) CreateMedium(url string) (path string, err error) {
	mock.CreateMockCalled = true
	fmt.Println("CHECK!!", url, mock.CreateMockCalled)
	return "path", nil
}
