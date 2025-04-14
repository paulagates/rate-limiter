package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %v", err)
	}
	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		err := r.client.Expire(ctx, key, expiration).Err()
		if err != nil {
			log.Println("Erro ao definir expiração:", err)
		}
	}
	return int(val), nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (int, error) {
	valStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RedisStorage) SetBlock(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key, "BLOCKED", duration).Err()
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "BLOCKED", nil
}

func (r *RedisStorage) ResetKey(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
