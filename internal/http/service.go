package http

import (
	"fmt"
	"net/http"
	"os"
	"rate-limiter/internal/config"
	"strconv"
	"time"

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

	handler := NewHandler()

	server := &http.Server{
		Addr:           s.cfg.Host + ":" + strconv.Itoa(s.cfg.Port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}
