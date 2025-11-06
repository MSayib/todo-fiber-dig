package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/msayib/todo-fiber-dig/internal/config"
)

func NewRedisClient(cfg config.Config) (*redis.Client, error) {
	redisAddress := fmt.Sprintf("%s:%s", cfg.RedisAddr, cfg.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Username: cfg.RedisUser,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	log.Println("Koneksi ke Redis (dengan username) berhasil!")
	return rdb, nil
}