package util

import (
	"context"
	"fmt"
	"golang-clean-architecture/internal/model"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterUtil struct {
	Redis      *redis.Client
	MaxRequest int64
	Duration   time.Duration
}

func NewRateLimiterUtil(redis *redis.Client) *RateLimiterUtil {
	maxReq, _ := strconv.ParseInt(os.Getenv("RATE_LIMIT_MAX_REQUEST"), 10, 64)
	if maxReq <= 0 {
		maxReq = 1
	}
	dur, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_DURATION_SEC"))
	if dur <= 0 {
		dur = 1
	}
	return &RateLimiterUtil{
		Redis:      redis,
		MaxRequest: maxReq,
		Duration:   time.Duration(dur) * time.Second,
	}
}

func (u RateLimiterUtil) IsAllowed(ctx context.Context, auth *model.Auth) bool {
	key := auth.ID.String()

	increment, err := u.Redis.Incr(ctx, key).Result()
	if err != nil {
		fmt.Println("Error incrementing:", err)
		return false
	}

	if increment == 1 {
		err := u.Redis.Expire(ctx, key, u.Duration).Err()
		if err != nil {
			fmt.Println("Error setting expiration:", err)
			return false
		}
	}

	return increment <= u.MaxRequest
}
