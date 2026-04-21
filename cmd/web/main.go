package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/sagemyrage/code-quality-expert-system/internal/config"
	apphttp "github.com/sagemyrage/code-quality-expert-system/internal/http"
)

func newPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.Postgres.DSN()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool error: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres error: %w", err)
	}

	return pool, nil
}

func newRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			log.Printf("close redis error: %v", closeErr)
		}
		return nil, fmt.Errorf("ping redis error: %w", err)
	}

	return client, nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	pgPool, err := newPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("init postgres error: %v", err)
	}
	defer pgPool.Close()

	redisClient, err := newRedisClient(ctx, cfg)
	if err != nil {
		log.Fatalf("init redis error: %v", err)
	}
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			log.Printf("close redis error: %v", closeErr)
		}
	}()

	router := apphttp.NewRouter()
	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}
	log.Printf("starting server on port %s", cfg.App.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("start server error: %v", err)
	}
}
