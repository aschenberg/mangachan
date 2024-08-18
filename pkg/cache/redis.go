package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"manga/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
		// DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		// ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		// WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		// PoolSize:     cfg.Redis.PoolSize,
		// PoolTimeout:  cfg.Redis.PoolTimeout,
	})

	res, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return &redis.Client{}, err
	}

	fmt.Printf(res)

	return redisClient, nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}

func Set[T any](ctx context.Context, c *redis.Client, key string, value T, duration time.Duration) error {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ct, key, v, duration).Err()
}

func Get[T any](ctx context.Context, c *redis.Client, key string) (T, error) {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()
	var dest T = *new(T)
	v, err := c.Get(ct, key).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func Del[T any](ctx context.Context, c *redis.Client, key string) error {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()

	return c.Del(ct, key).Err()
}
