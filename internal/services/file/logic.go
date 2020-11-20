package file

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"io"
	"mime"
	"net/http"
	"time"
)

type service struct {
	mc         *minio.Client
	bucketName string
}

func (svc *service) UploadFile(ctx context.Context, r io.Reader) (*UploadFileResponse, error) {
	buf := new(bytes.Buffer)

	fileSize, err := buf.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on read from input")
	}

	contentType := http.DetectContentType(buf.Bytes())

	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on get extensions by type")
	}

	ext := ""

	if len(exts) > 0 {
		ext = exts[0]
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	_, err = svc.mc.PutObject(ctx, svc.bucketName, fileName, buf, fileSize, minio.PutObjectOptions{
		ContentType:    contentType,
		SendContentMd5: false,
	})
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on put object")
	}

	rsp := UploadFileResponse{
		FileName: fileName,
	}
	return &rsp, nil
}

func (svc *service) DownloadFile(ctx context.Context, fileName string) (io.ReadSeeker, error) {
	res, err := svc.mc.GetObject(ctx, svc.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on get object")
	}

	return res, nil
}

func (svc *service) GetFileLastModified(ctx context.Context, fileName string) (*time.Time, error) {
	res, err := svc.mc.StatObject(ctx, svc.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), "error on get object")
	}

	return &res.LastModified, nil
}

func NewService(mc *minio.Client, bucketName string) Service {
	svc := service{
		mc:         mc,
		bucketName: bucketName,
	}

	return &svc
}
