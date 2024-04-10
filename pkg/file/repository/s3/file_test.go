package s3

import (
	"bytes"
	"doc-translate-go/mocks"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/mock/gomock"
)

func TestS3FileRepository_Persist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedUploader := mocks.NewMockUploaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"
	dat := []byte{}

	repo := NewS3FileRepository(mockedUploader, nil, nil, bucket)

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

func TestS3FileRepository_Persist_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedUploader := mocks.NewMockUploaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"
	dat := []byte{}

	repo := NewS3FileRepository(mockedUploader, nil, nil, bucket)

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

func TestS3FileRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDownloader := mocks.NewMockDownloaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"
	want := []byte("abc")

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	}

	repo := NewS3FileRepository(nil, mockedDownloader, nil, bucket)

	mockedDownloader.
		EXPECT().
		Download(gomock.Any(), gomock.Eq(input)).
		Do(func(w io.WriterAt, in *s3.GetObjectInput, args ...func(*s3manager.Downloader)) { w.WriteAt(want, 0) }).
		Times(1).
		Return(int64(1), nil)

	got, err := repo.Get(filepath)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestS3FileRepository_Get_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDownloader := mocks.NewMockDownloaderAPI(ctrl)

	bucket := "bucket"
	filepath := "file/path"

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	}

	repo := NewS3FileRepository(nil, mockedDownloader, nil, bucket)

	e := errors.New("get error")
	mockedDownloader.
		EXPECT().
		Download(gomock.Any(), gomock.Eq(input)).
		Times(1).
		Return(int64(0), e)

	got, err := repo.Get(filepath)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestS3FileRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mocks.NewMockS3API(ctrl)

	bucket := "bucket"
	filepath := "file/path"

	repo := NewS3FileRepository(nil, nil, mockedService, bucket)

	mockedService.EXPECT().DeleteObject(gomock.Eq(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	})).Return(nil, nil)

	mockedService.EXPECT().WaitUntilObjectNotExists(gomock.Eq(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	})).Return(nil)

	err := repo.Delete(filepath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3FileRepository_Delete_DeleteObjectErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mocks.NewMockS3API(ctrl)

	bucket := "bucket"
	filepath := "file/path"

	repo := NewS3FileRepository(nil, nil, mockedService, bucket)

	e := errors.New("delete error")
	mockedService.EXPECT().DeleteObject(gomock.Eq(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	})).Return(nil, e)

	err := repo.Delete(filepath)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestS3FileRepository_Delete_WaitErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mocks.NewMockS3API(ctrl)

	bucket := "bucket"
	filepath := "file/path"

	repo := NewS3FileRepository(nil, nil, mockedService, bucket)

	mockedService.EXPECT().DeleteObject(gomock.Eq(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	})).Return(nil, nil)

	e := errors.New("wait error")
	mockedService.EXPECT().WaitUntilObjectNotExists(gomock.Eq(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath),
	})).Return(e)

	err := repo.Delete(filepath)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}

func TestS3FileRepository_DeleteMany(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mocks.NewMockS3API(ctrl)

	bucket := "bucket"
	filepaths := []string{"file/path", "file/path2"}

	var objects []*s3.ObjectIdentifier
	for _, key := range filepaths {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(key)})
	}

	repo := NewS3FileRepository(nil, nil, mockedService, bucket)

	mockedService.EXPECT().DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}).Return(nil, nil)

	err := repo.DeleteMany(filepaths)
	if err != nil {
		t.Fatal(err)
	}
}

func TestS3FileRepository_DeleteMany_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mocks.NewMockS3API(ctrl)

	bucket := "bucket"
	filepaths := []string{"file/path", "file/path2"}

	repo := NewS3FileRepository(nil, nil, mockedService, bucket)

	e := errors.New("delete error")
	mockedService.EXPECT().DeleteObjects(gomock.Any()).Return(nil, e)

	err := repo.DeleteMany(filepaths)
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
}
