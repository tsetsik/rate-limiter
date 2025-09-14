package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tsetsik/rate-limiter/internal/config"

	"github.com/joho/godotenv"
)

type Service struct {
	cfg config.Config
}

func NewService() (*Service, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err.Error())
	}

	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	cfg, err := config.LoadConfig(host, port)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Service{
		cfg: *cfg,
	}, nil
}

func (s *Service) Start() error {
	// Validate the config
	if err := s.cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	mux := http.NewServeMux()

	httpResolver := NewHttpResolver()

	mux.Handle("GET /users", rateLimitMiddleware(http.HandlerFunc(httpResolver.GetUsers)))
	return http.ListenAndServe(s.cfg.Host+":"+strconv.Itoa(s.cfg.Port), mux)
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	rateLimiter := NewRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if !rateLimiter.Allow(userID) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
