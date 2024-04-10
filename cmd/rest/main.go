package main

import (
	"crypto/tls"
	"database/sql"
	"doc-translate-go/pkg/file/queue"
	filePG "doc-translate-go/pkg/file/repository/postgresql"
	fileS3 "doc-translate-go/pkg/file/repository/s3"
	fileUC "doc-translate-go/pkg/file/usecase"
	"doc-translate-go/pkg/tracker"
	"doc-translate-go/pkg/translator"
	userPG "doc-translate-go/pkg/user/repository/postgresql"
	userUC "doc-translate-go/pkg/user/usecase"
	documentprotov1 "doc-translate-go/proto/gen/go/proto/documentprocessor/v1"
	"doc-translate-go/rest/v1/handler"
	myMiddleware "doc-translate-go/rest/v1/middleware"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ENV_ADDR = "ADDR"

	ENV_DB_PASS = "DB_PASSWORD"
	ENV_DB_USER = "DB_USERNAME"
	ENV_DB_HOST = "DB_HOST"
	ENV_DB_PORT = "DB_PORT"
	ENV_DB_NAME = "DB"

	ENV_S3_BUCKET_NAME = "S3_BUCKET_NAME"
	ENV_AWS_REGION     = "AWS_REGION"

	ENV_TRANSLATE_GRPC_SERVER = "TRANSLATE_GRPC_SERVER"

	ENV_REDIS_ADDRS              = "REDIS_ADDRS"
	ENV_REDIS_PASS               = "REDIS_AUTH_TOKEN"
	ENV_REDIS_EXPIRATION_SECONDS = "REDIS_EXPIRATION_SECONDS"
)

var (
	db *sql.DB

	userUseCase           *userUC.UserUseCase
	fileUseCase           *fileUC.FileUseCase
	origFileMetaUseCase   *fileUC.OriginalFileMetadataUseCase
	translFileMetaUseCase *fileUC.TranslatedFileMetadataUseCase
	translateUseCase      *fileUC.TranslateUseCase
)

func init() {
	initDb()
	initUseCases()
}

func initDb() {
	var err error

	pass := url.QueryEscape(os.Getenv(ENV_DB_PASS))
	username := os.Getenv(ENV_DB_USER)
	host := os.Getenv(ENV_DB_HOST)
	port := os.Getenv(ENV_DB_PORT)
	dbname := os.Getenv(ENV_DB_NAME)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, pass, host, port, dbname)

	db, err = sql.Open("postgres", uri)
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
}

func initUseCases() {
	// Original File Metadata
	origFileMetaRepo := filePG.NewPostgresqlOriginalFileMetadataRepository(db)
	origFileMetaUseCase = fileUC.NewOriginalFileMetadataUseCase(origFileMetaRepo)

	// Translated File Metadata
	translFileMetaRepo := filePG.NewPostgresqlTranslatedFileMetadataRepository(db)
	translFileMetaUseCase = fileUC.NewTranslatedFileMetadataUseCase(translFileMetaRepo)

	// File
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv(ENV_AWS_REGION))})
	if err != nil {
		log.Fatalf("unable to get aws session: %v", err)
	}
	s3Uploader := s3manager.NewUploader(awsSession)
	s3Downloader := s3manager.NewDownloader(awsSession)
	fileRepo := fileS3.NewS3FileRepository(s3Uploader, s3Downloader, os.Getenv(ENV_S3_BUCKET_NAME))
	fileUseCase = fileUC.NewFileUseCase(fileRepo)

	// User
	userRepo := userPG.NewPostgresqlUserRepository(db)
	userUseCase = userUC.New(userRepo)

	// File Tracker
	grpcConn, err := grpc.Dial(os.Getenv(ENV_TRANSLATE_GRPC_SERVER), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to dial translate grpc server: %v", err)
	}
	grpcClient := documentprotov1.NewDocumentProcessorClient(grpcConn)
	grpcTranslator := translator.NewGrpcTranslator(grpcClient)

	redisAddrs := strings.Split(os.Getenv(ENV_REDIS_ADDRS), ",")
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:     redisAddrs,
		Password:  os.Getenv(ENV_REDIS_PASS),
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	expirationSecondsStr := os.Getenv(ENV_REDIS_EXPIRATION_SECONDS)
	expirationSeconds, err := strconv.Atoi(expirationSecondsStr)
	if err != nil {
		log.Fatalf("failed to parse redis expiration: %v", err)
	}
	redisFileTracker := tracker.NewRedisFileTracker(redisClient, expirationSeconds)

	err = redisFileTracker.Clear()
	if err != nil {
		log.Fatalf("failed to clear file tracker: %v", err)
	}

	// Translate Queue
	c := make(chan *queue.TranslateTask, 1<<32)
	chanTranslateQueue := queue.NewChannelTranslateQueue(c)

	// Translate
	translateUseCase = fileUC.NewTranslateUseCase(
		grpcTranslator,
		origFileMetaUseCase,
		translFileMetaUseCase,
		fileUseCase,
		redisFileTracker,
		chanTranslateQueue,
	)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: Middleware
	e.POST("/translate-docx", func(c echo.Context) error { return handler.TranslateDocx(c, translateUseCase) }, myMiddleware.AuthMiddleware(userUseCase))

	go translateUseCase.ListenAndExecute()

	e.Start(os.Getenv(ENV_ADDR))
}
