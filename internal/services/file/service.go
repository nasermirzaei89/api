package file

import (
	"context"
	"io"
)

type Service interface {
	UploadFile(ctx context.Context, r io.Reader) (res *UploadFileResponse, err error)
	DownloadFile(ctx context.Context, filename string) (res io.ReadSeeker, err error)
}

type UploadFileResponse struct {
	FileName string
}
