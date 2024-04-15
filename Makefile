protoc:
	@protoc --go_out=gen/go --go_opt=paths=source_relative --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative proto/documentprocessor/documentprocessor.proto

swag:
	@swagger init -g cmd/rest/main.go

mock:
	@mockgen -destination=mocks/mock_SQSAPI.go -package=mocks github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI
	@mockgen -destination=mocks/mock_UploaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI
	@mockgen -destination=mocks/mock_DownloaderAPI.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface DownloaderAPI
	@mockgen -destination=mocks/mock_S3API.go -package=mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API
	@mockgen -destination=mocks/mock_DocumentProcessClient.go -package=mocks doc-translate-go/gen/go/proto/documentprocessor DocumentProcessorClient

unittest:
	@go test -count=1 -coverprofile=coverage.out -short -race ./pkg/... ./rest/...
	@cat coverage.out | awk 'BEGIN {cov=0; stat=0;} $$3!="" { cov+=($$3==1?$$2:0); stat+=$$2; } END \
		{printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'

read-cov:
	@go tool cover -html=coverage.out

local-up:
	@docker compose -f docker-compose.local.yml up -d
	@migrate -source "file://migrations" -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

local-down:
	@docker compose -f docker-compose.local.yml down -v --remove-orphans

local-tls:
	@rm -rf tls/local/redis
	@mkdir -p tls/local/redis && cd tls/local/redis && curl -s https://raw.githubusercontent.com/redis/redis/cc0091f0f9fe321948c544911b3ea71837cf86e3/utils/gen-test-certs.sh | sh
	@mv tls/local/redis/tests/tls/* tls/local/redis/ && rm -r tls/local/redis/tests
