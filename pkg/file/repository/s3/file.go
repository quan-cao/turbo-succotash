package s3

import (
	"bytes"
	"doc-translate-go/pkg/file/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type S3FileRepository struct {
	uploader   s3manageriface.UploaderAPI
	downloader s3manageriface.DownloaderAPI
	service    s3iface.S3API
	bucketName string
}

func NewS3FileRepository(uploader s3manageriface.UploaderAPI, downloader s3manageriface.DownloaderAPI, service s3iface.S3API, bucketName string) *S3FileRepository {
	return &S3FileRepository{uploader, downloader, service, bucketName}
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

func (r *S3FileRepository) Get(filepath string) ([]byte, error) {
	buffer := aws.WriteAtBuffer{}

	_, err := r.downloader.Download(&buffer, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (r *S3FileRepository) Delete(filepath string) error {
	_, err := r.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return err
	}

	err = r.service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *S3FileRepository) DeleteMany(filepaths []string) error {
	var objects []*s3.ObjectIdentifier

	for _, key := range filepaths {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := r.service.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: &r.bucketName,
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Ensure implementation
var _ repository.FileRepository = (*S3FileRepository)(nil)
