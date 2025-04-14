package limiter

import (
	"net/http"
	"strings"

	"github.com/paulagates/rate-limiter/config"
)

func RateLimiterMiddleware(limiter *Limiter, config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("API_KEY")
			var identifier string
			var isToken bool

			if token != "" {
				identifier = token
				isToken = true
			} else {
				ip := strings.Split(r.RemoteAddr, ":")[0]
				identifier = ip
				isToken = false
			}

			allowed, err := limiter.AllowRequest(r.Context(), identifier, isToken)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
