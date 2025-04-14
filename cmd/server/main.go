package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/paulagates/rate-limiter/config"
	"github.com/paulagates/rate-limiter/internal/limiter"
	"github.com/paulagates/rate-limiter/internal/storage"
)

func main() {
	cfg := config.Load()
	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	rateLimiter := limiter.NewLimiter(redisStorage, cfg)

	r := chi.NewRouter()
	r.Use(limiter.RateLimiterMiddleware(rateLimiter, cfg))
	r.Get("/go", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "In Go We Trust!")
	})

	fmt.Println("Server Listening port: " + cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), r))
}
