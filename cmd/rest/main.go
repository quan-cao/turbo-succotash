package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"doc-translate-go/docs"
	"doc-translate-go/gen/go/proto/documentprocessor"
	"doc-translate-go/pkg/config"
	"doc-translate-go/pkg/file/queue"
	"doc-translate-go/pkg/tracker"
	"doc-translate-go/pkg/translator"
	"doc-translate-go/rest/v1/handler"
	"fmt"
	"log"

	filePG "doc-translate-go/pkg/file/repository/postgresql"
	fileS3 "doc-translate-go/pkg/file/repository/s3"
	fileUC "doc-translate-go/pkg/file/usecase"

	userPG "doc-translate-go/pkg/user/repository/postgresql"
	userUC "doc-translate-go/pkg/user/usecase"

	myMiddleware "doc-translate-go/rest/v1/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conf = config.NewConfig()

	db *sql.DB

	userUseCase           *userUC.UserUseCase
	authUseCase           *userUC.AuthUseCase
	fileUseCase           *fileUC.FileUseCase
	origFileMetaUseCase   *fileUC.OriginalFileMetadataUseCase
	translFileMetaUseCase *fileUC.TranslatedFileMetadataUseCase
	translateUseCase      *fileUC.TranslateUseCase
	progressUseCase       *fileUC.ProgressUseCase
)

func init() {
	initDb()
	initUseCases()
	initSwagger()
}

// @title DocsTranslateBackend
// @version 1.0
// @description API Routes for DocsTranslateBackend
// @BasePath /
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	addRoutes(e)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go translateUseCase.ListenAndExecute(ctx)

	err := e.Start(conf.App.Addr)
	if err != nil {
		panic(err)
	}
}

func initSwagger() {
	docs.SwaggerInfo.Host = conf.Swagger.Host
}

func initDb() {
	var err error

	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", conf.Db.Username, conf.Db.Password, conf.Db.Host, conf.Db.Port, conf.Db.DbName)

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
	awsSession, err := session.NewSession(&aws.Config{
		Region:   aws.String(conf.Aws.Region),
		Endpoint: aws.String(conf.Aws.Endpoint),
	})
	if err != nil {
		log.Fatalf("unable to get aws session: %v", err)
	}
	s3Uploader := s3manager.NewUploader(awsSession)
	s3Downloader := s3manager.NewDownloader(awsSession)
	s3Service := s3.New(awsSession)

	fileRepo := fileS3.NewS3FileRepository(s3Uploader, s3Downloader, s3Service, conf.Aws.S3BucketName)
	fileUseCase = fileUC.NewFileUseCase(fileRepo)

	// User
	userRepo := userPG.NewPostgresqlUserRepository(db)
	userUseCase = userUC.NewUserUseCase(userRepo)

	// Auth
	authUseCase = userUC.NewAuthUseCase(conf.Auth)

	// Translate
	fileTracker := getFileTracker()
	translateQueue := getTranslateQueue(awsSession)
	translr := getTranslator()
	translateUseCase = fileUC.NewTranslateUseCase(
		translr,
		origFileMetaUseCase,
		translFileMetaUseCase,
		fileUseCase,
		fileTracker,
		translateQueue,
	)

	// Progress
	progressUseCase = fileUC.NewProgressUseCase(fileTracker)
}

func getTranslator() translator.Translator {
	var translr translator.Translator

	switch conf.App.Translator {
	case "echo":
		translr = translator.NewEchoTranslator()
	default:
		grpcConn, err := grpc.Dial(conf.Translate.GrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("unable to dial translate grpc server: %v", err)
		}
		grpcClient := documentprocessor.NewDocumentProcessorClient(grpcConn)
		translr = translator.NewGrpcTranslator(grpcClient)
	}

	return translr
}

func getTranslateQueue(awsSession *session.Session) queue.TranslateQueue {
	var translateQueue queue.TranslateQueue

	switch conf.App.TranslateQueue {
	case "chan":
		c := make(chan *queue.TranslateTask, 1<<32)
		translateQueue = queue.NewChannelTranslateQueue(c)
	default:
		sqsClient := sqs.New(awsSession)
		translateQueue = queue.NewSqsTranslateQueue(sqsClient, conf.Aws.SqsQueueUrl, conf.Aws.SqsGroupId)
	}

	return translateQueue
}

func getFileTracker() tracker.FileTracker {
	var fileTracker tracker.FileTracker

	switch conf.App.FileTracker {
	default:
		redisClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    conf.Redis.Addrs,
			Password: conf.Redis.Password,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})

		fileTracker = tracker.NewRedisFileTracker(redisClient, conf.Redis.ExpirySeconds)

		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("failed to ping redis: %v", err)
		}
	}

	return fileTracker
}

func addRoutes(e *echo.Echo) {
	e.POST(
		"/translate-docx",
		func(c echo.Context) error { return handler.TranslateDocx(c, translateUseCase) },
		myMiddleware.AuthMiddleware(userUseCase, authUseCase),
	)

	e.DELETE(
		"/delete-translated-files",
		func(c echo.Context) error {
			return handler.DeleteFiles(c, userUseCase, origFileMetaUseCase, translFileMetaUseCase, fileUseCase)
		},
		myMiddleware.AuthMiddleware(userUseCase, authUseCase),
	)

	e.POST(
		"/download-translated-files",
		func(c echo.Context) error {
			return handler.DownloadTranslatedFiles(c, translFileMetaUseCase, fileUseCase)
		},
		myMiddleware.AuthMiddleware(userUseCase, authUseCase),
	)

	e.GET(
		"/show-translated-files",
		func(c echo.Context) error {
			return handler.ShowTranslatedFiles(c, translFileMetaUseCase)
		},
		myMiddleware.AuthMiddleware(userUseCase, authUseCase),
	)

	e.GET(
		"/upload-progress",
		func(c echo.Context) error {
			return handler.UploadProgress(c, progressUseCase)
		},
	)

	e.GET("/authorize", func(c echo.Context) error { return handler.Authorize(c, authUseCase) })

	e.GET("/token", func(c echo.Context) error { return handler.Token(c, authUseCase) })
}
