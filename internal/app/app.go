package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/nhassl3/hairdress_arz/internal/config"
	repoRedis "github.com/nhassl3/hairdress_arz/internal/repository/redis"
	transportGRPC "github.com/nhassl3/hairdress_arz/internal/transport/grpc"
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

	// migrations
	if cfg.Environment == "local" {
		if err := runMigrations(dsn, log); err != nil {
			return fmt.Errorf("app.Run: run migrations error: %w", err)
		}
	}

	// sqlc store initialize
	// TODO: initialize db store in internal
	//store := db.NewStore(pool)

	// redis
	redisSmsVerification, err := redisPkg.New(ctx, cfg.Redis.Address, cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("app.Run: init redis client error: %w", err)
	}
	defer func() { _ = redisSmsVerification.Close() }()
	log.Info("connected to redis (sms verification)")

	redisTokenBlackList, err := redisPkg.New(ctx, cfg.Redis.Address, cfg.Redis.Username, cfg.Redis.Password, cfg.Redis.DB+1)
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

	_ = accessManager

	refreshManager, err := auth.NewPasetoMaker(cfg.Auth.PasetoKey, cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return fmt.Errorf("app.Run: init auth error: %w", err)
	}

	_ = refreshManager

	// Register repositories

	// Register services
	svcs := &transportGRPC.Services{}

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

	if err := m.Up(); err != nil {
		return fmt.Errorf("app.runMigrations: up migrations error: %w", err)
	}

	log.Info("migration applied successfully")

	return nil
}
