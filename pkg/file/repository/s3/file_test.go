package s3

import (
	"bytes"
	"doc-translate-go/mocks"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/mock/gomock"
)

func TestS3FileRepository__Persist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedUploader := mocks.NewMockUploaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"
	dat := []byte{}

	repo := NewS3FileRepository(mockedUploader, nil, bucket)

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
		Body:   bytes.NewReader(dat),
	}

	mockedUploader.EXPECT().Upload(gomock.Eq(input)).Times(1)

	err := repo.Persist(dat, filepath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3FileRepository__Persist_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedUploader := mocks.NewMockUploaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"
	dat := []byte{}

	repo := NewS3FileRepository(mockedUploader, nil, bucket)

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
		Body:   bytes.NewReader(dat),
	}

	e := errors.New("upload error")
	mockedUploader.EXPECT().Upload(gomock.Eq(input)).Return(nil, e).Times(1)

	err := repo.Persist(dat, filepath)
	if err != e {
		t.Fatalf("expected error %v, got error %v", e, err)
	}
}
