package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	DB          DBConfig
	Redis       RedisConfig
	Auth        AuthConfig
	Log         LogConfig
	MinIO       MinIOConfig
}

type ServerConfig struct {
	GRPCPort,
	HTTPPort string
}

type DBConfig struct {
	Host  string
	Port  int
	User, // secret - from .env
	Password, // secret - from .env
	Name, // secret - from .env
	SSLMode string
}

type RedisConfig struct {
	Address,
	Username, // secret - from .env
	Password string // secret - from .env
	DB  int
	TTL RedisTTL
}

type RedisTTL struct {
	SmsVerificationCodeTTL,
	ProfileTTL,
	AuthBlockTTL time.Duration
}

type AuthConfig struct {
	PasetoKey string // secret - from .env
	AccessTokenTTL,
	RefreshTokenTTL time.Duration
}

type LogConfig struct {
	Level string
}

type MinIOConfig struct {
	Endpoint,
	AccessKey,
	SecretKey,
	Bucket string
	UseSSL bool
}

// LoadConfig reads public configuration from a YAML file and secrets from an .env file
//
// configFile - path to the YAML file
// envFile - path to the secrets .env file
func LoadConfig(configFile, envFile string) (*Config, error) {
	// ── YAML: public / non-sensitive settings ────────────────────────────────
	yv := viper.New()
	yv.SetConfigFile(configFile)
	yv.SetDefault("server.grpc_port", ":9090")
	yv.SetDefault("server.http_port", ":8080")
	yv.SetDefault("db.port", 5432)
	yv.SetDefault("db.ssl_mode", "disable")
	yv.SetDefault("auth.access_token_ttl", "15m")
	yv.SetDefault("auth.refresh_token_ttl", "168h")
	yv.SetDefault("redis.db", 0)
	yv.SetDefault("log.level", "info")
	yv.SetDefault("redis.ttl.sms_code_verification", "5m")
	yv.SetDefault("redis.ttl.profile", "15m")
	yv.SetDefault("redis.ttl.auth_block", "5m")
	yv.SetDefault("minio.endpoint", "localhost:9000")
	yv.SetDefault("minio.use_ssl", "true")

	if err := yv.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read yaml %q: %w", configFile, err)
	}

	yv.BindEnv("db.host", "DB_HOST")
	yv.BindEnv("db.port", "DB_PORT")
	yv.BindEnv("redis.addr", "REDIS_ADDR")
	yv.BindEnv("minio.endpoint", "MINIO_ENDPOINT")

	// ── .env: secrets ─────────────────────────────────────────────────────────
	ev := viper.New()
	ev.SetConfigFile(envFile)
	ev.SetConfigType("env")
	if err := ev.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read env %q: %w", envFile, err)
	}

	// ── Assemble ──────────────────────────────────────────────────────────────
	cfg := &Config{}
	cfg.Environment = yv.GetString("environment")
	cfg.Server.GRPCPort = yv.GetString("server.grpc_port")
	cfg.Server.HTTPPort = yv.GetString("server.http_port")

	cfg.DB.Host = yv.GetString("db.host")
	cfg.DB.Port = yv.GetInt("db.port")
	cfg.DB.SSLMode = yv.GetString("db.ssl_mode")
	cfg.DB.User = ev.GetString("DB_USER")
	cfg.DB.Password = ev.GetString("DB_PASSWORD")
	cfg.DB.Name = ev.GetString("DB_NAME")

	cfg.Redis.Address = yv.GetString("redis.addr")
	cfg.Redis.Username = ev.GetString("REDIS_USER")
	cfg.Redis.Password = ev.GetString("REDIS_USER_PASSWORD")
	cfg.Redis.DB = yv.GetInt("redis.db")
	cfg.Redis.TTL.SmsVerificationCodeTTL = yv.GetDuration("redis.ttl.sms_code_verification")
	cfg.Redis.TTL.ProfileTTL = yv.GetDuration("redis.ttl.profile")
	cfg.Redis.TTL.AuthBlockTTL = yv.GetDuration("redis.ttl.auth_block")
	cfg.Auth.AccessTokenTTL = yv.GetDuration("auth.access_token_ttl")
	cfg.Auth.RefreshTokenTTL = yv.GetDuration("auth.refresh_token_ttl")
	cfg.Auth.PasetoKey = ev.GetString("PASETO_KEY")

	cfg.MinIO.Endpoint = yv.GetString("minio.endpoint")
	cfg.MinIO.AccessKey = ev.GetString("MINIO_ACCESS_KEY")
	cfg.MinIO.SecretKey = ev.GetString("MINIO_SECRET_KEY")
	cfg.MinIO.Bucket = yv.GetString("minio.bucket")
	cfg.MinIO.UseSSL = yv.GetBool("minio.use_ssl")

	cfg.Log.Level = yv.GetString("log.level")

	return cfg, nil
}
