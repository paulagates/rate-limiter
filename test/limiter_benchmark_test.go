package limiter_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/paulagates/rate-limiter/config"
	"github.com/paulagates/rate-limiter/internal/limiter"
	"github.com/paulagates/rate-limiter/internal/storage"
)

func BenchmarkAllowRequest(b *testing.B) {
	cfg := config.Load()
	storage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Erro ao conectar com Redis: %v", err)
	}

	lim := limiter.NewLimiter(storage, cfg)

	ctx := context.Background()
	identifier := "127.0.0.1"

	storage.ResetKey(ctx, identifier)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := lim.AllowRequest(ctx, identifier, false)
		if err != nil {
			b.Errorf("Error allowing request: %v", err)
		}
	}
}

func BenchmarkLimiter_AllowRequest(b *testing.B) {
	ctx := context.Background()
	cfg := config.Load()
	store, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		b.Fatalf("Error connecting to Redis: %v", err)
	}

	lim := limiter.NewLimiter(store, cfg)

	identifier := "127.0.0.1"

	store.SetBlock(ctx, "block:"+identifier, time.Second*1)
	time.Sleep(time.Second * 2)

	b.ResetTimer()

	var wg sync.WaitGroup
	numGoroutines := 100

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 5 {
				allowed, err := lim.AllowRequest(ctx, identifier, false)
				if err != nil {
					b.Logf("erro: %v", err)
				} else {
					b.Logf("allowed=%v", allowed)
				}
			}
		}()
	}

	wg.Wait()
}
