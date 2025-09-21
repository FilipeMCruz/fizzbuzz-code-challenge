package main

import (
	"context"
	"fizzbuzz-code-challenge/handlers"
	"fizzbuzz-code-challenge/infrastructure"
	"flag"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	running := func() {
		slog.Info("fizzbuzz server running", "port", *port)
	}

	err := start(ctx, stop, running, *port)
	if err != nil {
		slog.Error("failed to start service", "error", err)
	}
}

// start registers the handlers (wrapped with logging and stats) in a ServeMux
// and calls infrastructure.Run to run the http Server.
func start(ctx context.Context, stop func(), running func(), port int) error {
	ch := make(chan string)
	defer close(ch)

	warp := infrastructure.BuildWrapHandlerChain(ch)

	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/fizzbuzz", warp(handlers.BuildFizzBuzzHandler()))
	mux.Handle("GET /api/v1/stats", warp(handlers.BuildStatsHandler(ch)))

	return infrastructure.Run(ctx, stop, running, port, mux)
}
