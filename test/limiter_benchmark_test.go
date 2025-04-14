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

func setupLimiter() (*limiter.Limiter, *storage.RedisStorage, context.Context) {
	cfg := config.Load()
	storage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Erro ao conectar com Redis: %v", err)
	}

	return limiter.NewLimiter(storage, cfg), storage, context.Background()
}
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
			b.Errorf("Erro ao permitir requisição: %v", err)
		}
	}
}

func BenchmarkLimiter_IP(b *testing.B) {
	lim, store, ctx := setupLimiter()
	ip := "127.0.0.1"

	_ = store.ResetKey(ctx, ip)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := lim.AllowRequest(ctx, ip, false)
		if err != nil {
			b.Errorf("Erro na requisição IP: %v", err)
		}
	}
}

func BenchmarkLimiter_Token(b *testing.B) {
	lim, store, ctx := setupLimiter()
	token := "token:abc123"

	_ = store.ResetKey(ctx, token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := lim.AllowRequest(ctx, token, true)
		if err != nil {
			b.Errorf("Erro na requisição Token: %v", err)
		}
	}
}

func BenchmarkLimiter_AllowRequest(b *testing.B) {
	ctx := context.Background()
	cfg := config.Load()

	store, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		b.Fatalf("erro ao conectar ao Redis: %v", err)
	}

	lim := limiter.NewLimiter(store, cfg)

	identifier := "127.0.0.1"

	store.SetBlock(ctx, "block:"+identifier, time.Second*1)
	time.Sleep(time.Second * 2)

	b.ResetTimer()

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
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
