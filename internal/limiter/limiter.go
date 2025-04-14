package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/paulagates/rate-limiter/config"
	"github.com/paulagates/rate-limiter/internal/storage"
)

type Limiter struct {
	storage storage.Storage
	cfg     *config.Config
}

func NewLimiter(s storage.Storage, cfg *config.Config) *Limiter {
	return &Limiter{
		storage: s,
		cfg:     cfg,
	}
}

func (l *Limiter) AllowRequest(ctx context.Context, identifier string, isToken bool) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", identifier)

	blocked, err := l.storage.IsBlocked(ctx, blockKey)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	var limit int
	var expire time.Duration
	if isToken {
		limit = l.cfg.RateLimitToken
		expire = time.Duration(l.cfg.BlockDurationToken) * time.Second
	} else {
		limit = l.cfg.RateLimitIP
		expire = time.Duration(l.cfg.BlockDurationIP) * time.Second
	}

	count, err := l.storage.Increment(ctx, identifier, time.Second)
	if err != nil {
		return false, err
	}
	if count > limit {
		err := l.storage.SetBlock(ctx, blockKey, expire)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
