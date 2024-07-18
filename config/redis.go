package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisC struct {
	client *redis.Client
	env    *Env
}

func NewKeyDB(env *Env) *redis.Client {
	Host := env.RedisHost
	Port := env.RedisPort
	Password := env.RedisPassword
	print(fmt.Sprintf("%s:%d", Host, Port))
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Host, Port),
		Password: Password,
		DB:       0,
	}
	return redis.NewClient(opts)
}

// NewTask instantiates the Task repository.
func NewRedis(client *redis.Client) *RedisC {
	return &RedisC{
		client: client,
	}
}

func (r *RedisC) Set(c context.Context, key string, value interface{}, expire time.Duration) error {
	var err error
	ctx, cancel := context.WithTimeout(c, r.env.RedisContextTime)
	defer cancel()
	enc, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, enc, expire).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *RedisC) Get(c context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(c, r.env.RedisContextTime)
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", key)
	}
	var email string

	err = json.Unmarshal([]byte(val), &email)

	return email, err
}

func (r *RedisC) Del(c context.Context, key string) error {
	ctx, cancel := context.WithTimeout(c, r.env.RedisContextTime)
	defer cancel()

	return r.client.Del(ctx, key).Err()
}
