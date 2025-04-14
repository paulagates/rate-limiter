package storage

import (
	"context"
	"time"
)

type Storage interface {
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

	Get(ctx context.Context, key string) (int, error)

	SetBlock(ctx context.Context, key string, duration time.Duration) error

	IsBlocked(ctx context.Context, key string) (bool, error)

	ResetKey(ctx context.Context, key string) error
}
