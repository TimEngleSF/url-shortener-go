package S3

import "context"

type S3Mock struct{}

func (s3 *S3Mock) UploadFile(ctx context.Context, key string, img []byte) error {
	return nil
}
