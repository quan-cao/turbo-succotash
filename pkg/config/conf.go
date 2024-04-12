package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	ENV_ADDR = "ADDR"

	ENV_DB_USERNAME = "DB_USERNAME"
	ENV_DB_PASSWORD = "DB_PASSWORD"
	ENV_DB_HOST     = "DB_HOST"
	ENV_DB_PORT     = "DB_PORT"
	ENV_DB_NAME     = "DB"

	ENV_AWS_S3_BUCKET_NAME = "S3_BUCKET_NAME"
	ENV_AWS_REGION         = "AWS_REGION"

	ENV_TRANSLATE_GRPC_SERVER = "TRANSLATE_GRPC_SERVER"

	ENV_REDIS_ADDRS          = "REDIS_ADDRS"
	ENV_REDIS_PASSWORD       = "REDIS_AUTH_TOKEN"
	ENV_REDIS_EXPIRY_SECONDS = "REDIS_EXPIRATION_SECONDS"

	ENV_AUTH_TOKEN_ENDPOINT             = "AUTH_TOKEN_ENDPOINT"
	ENV_AUTH_ENDPOINT                   = "AUTH_ENDPOINT"
	ENV_AUTH_INTROSPECT_ENDPOINT        = "AUTH_INTROSPECT_ENDPOINT"
	ENV_AUTH_USERINFO_ENDPOINT          = "AUTH_USERINFO_ENDPOINT"
	ENV_AUTH_CLIENT_ID                  = "AUTH_CLIENT_ID"
	ENV_AUTH_CLIENT_SECRET              = "AUTH_CLIENT_SECRET"
	ENV_AUTH_DISTRIBUTION_LIST_ENDPOINT = "AUTH_DISTRIBUTION_LIST_ENDPOINT"
	ENV_AUTH_DISTRIBUTION_LIST          = "AUTH_DISTRIBUTION_LIST"
)

type Config struct {
	App       *AppConfig
	Aws       *AwsConfig
	Translate *TranslateConfig
	Redis     *RedisConfig
	Auth      *AuthConfig
	Db        *DbConfig
}

func NewConfig() *Config {
	return &Config{
		App:       NewAppConfig(),
		Aws:       NewAwsConfig(),
		Translate: NewTranslateConfig(),
		Redis:     NewRedisConfig(),
		Auth:      NewAuthConfig(),
		Db:        NewDbconfig(),
	}
}

type AppConfig struct {
	Addr string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Addr: os.Getenv(ENV_ADDR),
	}
}

type AwsConfig struct {
	Region       string
	S3BucketName string
}

func NewAwsConfig() *AwsConfig {
	return &AwsConfig{
		Region:       os.Getenv(ENV_AWS_S3_BUCKET_NAME),
		S3BucketName: os.Getenv(ENV_AWS_S3_BUCKET_NAME),
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
