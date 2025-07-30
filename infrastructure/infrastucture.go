package infrastructure

import (
	"context"
	"errors"
	"fizzbuzz-code-challenge/infrastructure/logging"
	"fizzbuzz-code-challenge/infrastructure/stats"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// BuildWrapHandlerChain creates a function that wraps a handler with:
// - request counter (for stats)
// - basic request info logging
func BuildWrapHandlerChain(ch chan<- string) func(http.Handler) http.Handler {
	handler := stats.BuildWrapStats(ch)

	return func(next http.Handler) http.Handler {
		next = handler(next)
		return logging.WrapLogging(next)
	}
}

// Run runs an http server and ensures that it is gracefully shutdown:
// - in flight requests are answered
// - new requests are not accepted
func Run(ctx context.Context, stop func(), port int, handler http.Handler) error {
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
		
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Printf("fizzbuzz server starting on port %d", port)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer stopOngoingGracefully()

	return httpServer.Shutdown(shutdownCtx)
}
