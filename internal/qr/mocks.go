package qr

type QRCodeMock struct {
	CreateMockCalled bool
}

func (mock *QRCodeMock) CreateMedium(url string) ([]byte, error) {
	mock.CreateMockCalled = true
	return []byte{}, nil
}

func (mock *QRCodeMock) CreateMediumFile(url string) (path string, err error) {
	mock.CreateMockCalled = true
	return "path", nil
}
