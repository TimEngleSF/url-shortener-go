package qr

type QRCodeMock struct {
	CreateMockCalled bool
}

func (mock *QRCodeMock) CreateMedium(url string) (path string, err error) {
	mock.CreateMockCalled = true
	return "path", nil
}
