package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/sagemyrage/code-quality-expert-system/internal/config"
)

func newPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.Postgres.DSN()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
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
			log.Printf("close redis: %v", closeErr)
		}
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pgPool, err := newPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("init postgres: %v", err)
	}
	defer pgPool.Close()

	redisClient, err := newRedisClient(ctx, cfg)
	if err != nil {
		log.Fatalf("init redis: %v", err)
	}
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			log.Printf("close redis: %v", closeErr)
		}
	}()

	log.Printf("starting server on port %s", cfg.App.Port)
}
