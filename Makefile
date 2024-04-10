mock:
	@mockgen -destination=mocks/mock_SQSAPI.go -package=mocks github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI
	@mockgen -destination=mocks/mock_UploaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI
	@mockgen -destination=mocks/mock_DownloaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface DownloaderAPI
	@mockgen -destination=mocks/mock_S3API.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API
	@mockgen -destination=mocks/mock_DocumentProcessClient.go -package=mocks doc-translate-go/proto/gen/go/proto/documentprocessor/v1 DocumentProcessorClient

unit:
	@go test -coverprofile=coverage.out -short ./... ; \
		cat coverage.out | \
		awk 'BEGIN {cov=0; stat=0;} $$3!="" { cov+=($$3==1?$$2:0); stat+=$$2; } END \
		{printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'

covhtml:
	@go tool cover -html=coverage.out
