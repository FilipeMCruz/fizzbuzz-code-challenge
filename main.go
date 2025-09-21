package main

import (
	"context"
	"fizzbuzz-code-challenge/handlers"
	"fizzbuzz-code-challenge/infrastructure"
	"fizzbuzz-code-challenge/services"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	log.Printf("fizzbuzz server running on port: %d", port)
	err := start(ctx, stop, func() {}, *port)
	if err != nil {
		log.Printf("failed to start service: %v", err)
	}
}

// start registers the handlers (wrapped with logging, stats and recovery) in a ServeMux
// and calls infrastructure.Run to run the http Server.
func start(ctx context.Context, stop func(), running func(), port int) error {
	s := services.NewStats()

	warp := infrastructure.BuildWrapHandlerChain(s.Increment)

	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/fizzbuzz", warp(handlers.BuildFizzBuzzHandler()))
	mux.Handle("GET /api/v1/stats", warp(handlers.BuildStatsHandler(s)))

	return infrastructure.Run(ctx, stop, running, port, mux)
}
