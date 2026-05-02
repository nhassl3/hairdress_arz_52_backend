package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nhassl3/hairdress_arz/internal/config"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
	postgresRedis "github.com/nhassl3/hairdress_arz/internal/repository/postgres"
	repoRedis "github.com/nhassl3/hairdress_arz/internal/repository/redis"
	"github.com/nhassl3/hairdress_arz/internal/service"
	transportGRPC "github.com/nhassl3/hairdress_arz/internal/transport/grpc"
	"github.com/nhassl3/hairdress_arz/pkg/sms"
	"github.com/nhassl3/servicehub-backend/pkg/auth"
	"github.com/nhassl3/servicehub-backend/pkg/logger"
	"github.com/nhassl3/servicehub-backend/pkg/postgres"
	redisPkg "github.com/nhassl3/servicehub-backend/pkg/redis"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) error {
	// Initialize Logger (ZAP)
	log, err := logger.NewZapLogger(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("app.Run: init logger error: %w", err)
	}
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(log)

	// Database init
	ctx := context.Background()
	dsn := postgres.DSN(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)

	pool, err := postgres.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("app.Run: init postgres error: %w", err)
	}
	defer pool.Close()
	log.Info("connected to PostgresSQL")

	// Init migrations, sms sender object
	var smsSender domain.SmsSender
	if cfg.Environment == "local" {
		if err := runMigrations(dsn, log); err != nil {
			return fmt.Errorf("app.Run: run migrations error: %w", err)
		}
		smsSender = sms.NewSMSSenderLog(log, cfg.Auth.OTPConfig.CodeLength, cfg.Auth.OTPConfig.SecretKey)
	} else {
		smsSender = sms.NewSMSender(cfg.Auth.OTPConfig.CodeLength, cfg.Auth.OTPConfig.SecretKey)
	}

	// SQLC store initialize
	store := db.NewStore(pool)

	// Connect Redis
	userRedis, err := redisPkg.New(ctx, cfg.Redis.Address, cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("app.Run: init redis store error: %w", err)
	}
	defer func() { _ = userRedis.Close() }()
	log.Info("connected to redis (user)")

	userImplementsRedis := repoRedis.NewUserRedis(userRedis, cfg.Redis.TTL.ProfileTTL, cfg.Redis.TTL.AuthBlockTTL)

	smsVerificationRedis, err := redisPkg.New(ctx, cfg.Redis.Address, cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.DB+1)
	if err != nil {
		return fmt.Errorf("app.Run: init redis client error: %w", err)
	}
	defer func() { _ = smsVerificationRedis.Close() }()
	log.Info("connected to redis (sms verification)")

	smsVerificationImplementsRedis := repoRedis.NewSMSRedis(
		smsVerificationRedis,
		cfg.Redis.TTL.SmsVerificationCodeTTL,
		cfg.Auth.OTPConfig.Cooldown,
		cfg.Auth.OTPConfig.Attempts,
		cfg.Auth.OTPConfig.DailyPerPhone,
		cfg.Auth.OTPConfig.DailyPerIP,
	)

	redisTokenBlackList, err := redisPkg.New(ctx, cfg.Redis.Address, cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.DB+2)
	if err != nil {
		return fmt.Errorf("app.Run: init redis client error: %w", err)
	}
	defer func() { _ = redisTokenBlackList.Close() }()
	log.Info("connected to redis (token blacklist)")

	tokenBlacklist := repoRedis.NewTokenBlacklist(redisTokenBlackList)

	// MinIO initialize
	//minIOClient, err := minio.NewMinIO(
	//	ctx, cfg.MinIO.Endpoint, cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, "", cfg.MinIO.Bucket, cfg.MinIO.UseSSL,
	//)
	//if err != nil {
	//	return fmt.Errorf("app.Run: init minio client error: %w", err)
	//}
	//log.Info("connected to minio")

	// token managers initialize
	accessMaker, err := auth.NewPasetoMaker(cfg.Auth.PasetoKey, cfg.Auth.AccessTokenTTL)
	if err != nil {
		return fmt.Errorf("app.Run: init auth error: %w", err)
	}
	accessManager := auth.NewBlacklistedTokenManager(accessMaker, tokenBlacklist)

	refreshManager, err := auth.NewPasetoMaker(cfg.Auth.PasetoKey, cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return fmt.Errorf("app.Run: init auth error: %w", err)
	}

	// Register repositories
	authRepo := postgresRedis.NewAuthRepo(store)

	// Register services
	svcs := &transportGRPC.Services{
		Auth: service.NewAuthService(authRepo,
			accessManager,
			refreshManager,
			tokenBlacklist,
			userImplementsRedis,
			smsVerificationImplementsRedis,
			smsSender,
		),
	}

	// make gRPC server
	grpcServer := transportGRPC.NewServer(svcs, accessManager, log)

	errCh := make(chan error, 2)

	go func() {
		if err := grpcServer.Start(cfg.Server.GRPCPort); err != nil {
			errCh <- fmt.Errorf("app.Run: start grpc server error: %w", err)
		}
	}()

	go func() {
		if err := grpcServer.StartGateway(ctx, "localhost"+cfg.Server.GRPCPort, cfg.Server.HTTPPort); err != nil {
			errCh <- fmt.Errorf("app.Run: start grpc server error: %w", err)
		}
	}()

	log.Info("ArzSalon server started",
		zap.String("port", cfg.Server.GRPCPort),
		zap.String("http port", cfg.Server.HTTPPort),
		zap.String("env", cfg.Environment),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		log.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*1e9)
		defer cancel()
		return grpcServer.Shutdown(shutdownCtx)
	}
}

func runMigrations(dsn string, log *zap.Logger) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("app.runMigrations: init migrations error: %w", err)
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("app.runMigrations: up migrations error: %w", err)
	}

	log.Info("migration applied successfully")

	return nil
}
