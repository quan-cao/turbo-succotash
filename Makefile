mock:
	mockgen -destination=mocks/mock_SQSAPI.go -package=mocks github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI
	mockgen -destination=mocks/mock_UploaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI
	mockgen -destination=mocks/mock_DownloaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface DownloaderAPI
	mockgen -destination=mocks/mock_DocumentProcessClient.go -package=mocks doc-translate-go/proto/gen/go/proto/documentprocessor/v1 DocumentProcessorClient
	
