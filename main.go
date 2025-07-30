package main

import (
	"context"
	"fizzbuzz-code-challenge/handlers"
	"fizzbuzz-code-challenge/infrastructure"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	err := start(ctx, stop, *port)
	if err != nil {
		log.Fatal(err)
	}
}

// start registers the handlers (wrapped with logging and stats) in a ServeMux
// and calls infrastructure.Run to run the http Server
func start(ctx context.Context, stop func(), port int) error {
	ch := make(chan string)
	defer close(ch)

	warp := infrastructure.BuildWrapHandlerChain(ch)

	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/fizzbuzz", warp(handlers.BuildFizzBuzzHandler()))
	mux.Handle("GET /api/v1/stats", warp(handlers.BuildStatsHandler(ch)))

	return infrastructure.Run(ctx, stop, port, mux)
}
