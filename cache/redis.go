package cache

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func Initialize() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
		return err
	}
	log.Printf("connected to redis server: %v", redisAddr)
	return nil
}

func GetRedisClient() *redis.Client {
	return rdb
}

func GetRedisCtx() context.Context {
	return ctx
}
