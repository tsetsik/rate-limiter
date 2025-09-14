package main

import (
	service "github.com/tsetsik/rate-limiter/internal/http"
)

func main() {
	// This is a placeholder for the rate limiter implementation.
	// Actual implementation would go here.

	svc, err := service.NewService()
	if err != nil {
		panic(err)
	}

	if err := svc.Start(); err != nil {
		panic(err)
	}
}
