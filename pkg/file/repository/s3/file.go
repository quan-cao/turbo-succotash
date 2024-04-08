package s3

import (
	"bytes"
	"doc-translate-go/pkg/file/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type S3FileRepository struct {
	uploader   s3manageriface.UploaderAPI
	downloader s3manageriface.DownloaderAPI
	bucketName string
}

func NewS3FileRepository(uploader s3manageriface.UploaderAPI, downloader s3manageriface.DownloaderAPI, bucketName string) *S3FileRepository {
	return &S3FileRepository{uploader, downloader, bucketName}
}

func (r *S3FileRepository) Persist(b []byte, filepath string) error {
	_, err := r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filepath),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *S3FileRepository) Get(filepath string) ([]byte, error) { panic("") }

func (r *S3FileRepository) Delete(filepath string) error { panic("") }

func (r *S3FileRepository) DeleteMany(filepaths []string) error { panic("") }

// Ensure implementation
var _ repository.FileRepository = (*S3FileRepository)(nil)
