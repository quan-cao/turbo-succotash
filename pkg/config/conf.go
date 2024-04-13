package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	ENV_ADDR            = "ADDR"
	ENV_TRANSLATOR      = "TRANSLATOR"
	ENV_TRANSLATE_QUEUE = "TRANSLATE_QUEUE"
	ENV_FILE_TRACKER    = "FILE_TRACKER"

	ENV_DB_USERNAME = "DB_USERNAME"
	ENV_DB_PASSWORD = "DB_PASSWORD"
	ENV_DB_HOST     = "DB_HOST"
	ENV_DB_PORT     = "DB_PORT"
	ENV_DB_NAME     = "DB_NAME"

	ENV_AWS_REGION   = "AWS_REGION"
	ENV_AWS_ENDPOINT = "AWS_ENDPOINT"

	ENV_AWS_S3_BUCKET_NAME = "AWS_S3_BUCKET_NAME"

	ENV_AWS_SQS_QUEUE_URL = "AWS_SQS_QUEUE_URL"
	ENV_AWS_SQS_GROUP_ID  = "AWS_SQS_GROUP_ID"

	ENV_TRANSLATE_GRPC_SERVER = "TRANSLATE_GRPC_SERVER"

	ENV_REDIS_ADDRS          = "REDIS_ADDRS"
	ENV_REDIS_PASSWORD       = "REDIS_PASSWORD"
	ENV_REDIS_EXPIRY_SECONDS = "REDIS_EXPIRY_SECONDS"

	ENV_AUTH_TOKEN_ENDPOINT             = "AUTH_TOKEN_ENDPOINT"
	ENV_AUTH_ENDPOINT                   = "AUTH_ENDPOINT"
	ENV_AUTH_INTROSPECT_ENDPOINT        = "AUTH_INTROSPECT_ENDPOINT"
	ENV_AUTH_USERINFO_ENDPOINT          = "AUTH_USERINFO_ENDPOINT"
	ENV_AUTH_CLIENT_ID                  = "AUTH_CLIENT_ID"
	ENV_AUTH_CLIENT_SECRET              = "AUTH_CLIENT_SECRET"
	ENV_AUTH_DISTRIBUTION_LIST_ENDPOINT = "AUTH_DISTRIBUTION_LIST_ENDPOINT"
	ENV_AUTH_DISTRIBUTION_LIST          = "AUTH_DISTRIBUTION_LIST"

	ENV_SWAGGER_HOST = "SWAGGER_HOST"
)

type Config struct {
	App       *AppConfig
	Aws       *AwsConfig
	Translate *TranslateConfig
	Redis     *RedisConfig
	Auth      *AuthConfig
	Db        *DbConfig
	Swagger   *SwaggerConfig
}

func NewConfig() *Config {
	return &Config{
		App:       NewAppConfig(),
		Aws:       NewAwsConfig(),
		Translate: NewTranslateConfig(),
		Redis:     NewRedisConfig(),
		Auth:      NewAuthConfig(),
		Db:        NewDbconfig(),
		Swagger:   NewSwaggerConfig(),
	}
}

type AppConfig struct {
	Addr           string
	Translator     string
	FileTracker    string
	TranslateQueue string
}

func NewAppConfig() *AppConfig {
	addr := os.Getenv(ENV_ADDR)
	if addr == "" {
		addr = ":8080"
	}

	translator := os.Getenv(ENV_TRANSLATOR)
	if translator == "" {
		translator = "grpc"
	}

	translateQueue := os.Getenv(ENV_TRANSLATE_QUEUE)
	if translateQueue == "" {
		translateQueue = "sqs"
	}

	fileTracker := os.Getenv(ENV_FILE_TRACKER)
	if fileTracker == "" {
		fileTracker = "redis"
	}

	return &AppConfig{
		Addr:           addr,
		Translator:     translator,
		TranslateQueue: translateQueue,
		FileTracker:    fileTracker,
	}
}

type AwsConfig struct {
	Region       string
	Endpoint     string
	S3BucketName string
	SqsQueueUrl  string
	SqsGroupId   string
}

func NewAwsConfig() *AwsConfig {
	return &AwsConfig{
		Region:       os.Getenv(ENV_AWS_REGION),
		Endpoint:     os.Getenv(ENV_AWS_ENDPOINT),
		S3BucketName: os.Getenv(ENV_AWS_S3_BUCKET_NAME),
		SqsQueueUrl:  os.Getenv(ENV_AWS_SQS_QUEUE_URL),
		SqsGroupId:   os.Getenv(ENV_AWS_SQS_GROUP_ID),
	}
}

type TranslateConfig struct {
	GrpcServer string
}

func NewTranslateConfig() *TranslateConfig {
	return &TranslateConfig{
		GrpcServer: os.Getenv(ENV_TRANSLATE_GRPC_SERVER),
	}
}

type AuthConfig struct {
	AuthEndpoint             string
	TokenEndpoint            string
	IntrospectEndpoint       string
	UserInfoEndpoint         string
	ClientId                 string
	ClientSecret             string
	DistributionListEndpoint string
	DistributionList         []string
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		TokenEndpoint:            os.Getenv(ENV_AUTH_TOKEN_ENDPOINT),
		AuthEndpoint:             os.Getenv(ENV_AUTH_ENDPOINT),
		IntrospectEndpoint:       os.Getenv(ENV_AUTH_INTROSPECT_ENDPOINT),
		UserInfoEndpoint:         os.Getenv(ENV_AUTH_USERINFO_ENDPOINT),
		ClientId:                 os.Getenv(ENV_AUTH_CLIENT_ID),
		ClientSecret:             os.Getenv(ENV_AUTH_CLIENT_SECRET),
		DistributionListEndpoint: os.Getenv(ENV_AUTH_DISTRIBUTION_LIST_ENDPOINT),
		DistributionList:         strings.Split(strings.TrimSpace(os.Getenv(ENV_AUTH_DISTRIBUTION_LIST)), ","),
	}
}

type RedisConfig struct {
	Addrs         []string
	Password      string
	ExpirySeconds int
}

func NewRedisConfig() *RedisConfig {
	expiry, err := strconv.Atoi(os.Getenv(ENV_REDIS_EXPIRY_SECONDS))
	if err != nil {
		expiry = 0
	}

	return &RedisConfig{
		Addrs:         strings.Split(strings.TrimSpace(os.Getenv(ENV_REDIS_ADDRS)), ","),
		Password:      os.Getenv(ENV_REDIS_PASSWORD),
		ExpirySeconds: expiry,
	}
}

type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

func NewDbconfig() *DbConfig {
	port, err := strconv.Atoi(os.Getenv(ENV_DB_PORT))
	if err != nil {
		port = 0
	}

	return &DbConfig{
		Username: os.Getenv(ENV_DB_USERNAME),
		Password: os.Getenv(ENV_DB_PASSWORD),
		Host:     os.Getenv(ENV_DB_HOST),
		Port:     port,
		DbName:   os.Getenv(ENV_DB_NAME),
	}
}

type SwaggerConfig struct {
	Host string
}

func NewSwaggerConfig() *SwaggerConfig {
	return &SwaggerConfig{
		Host: os.Getenv(ENV_SWAGGER_HOST),
	}
}
