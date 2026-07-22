package config

import (
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

func NewRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return rdb
}
