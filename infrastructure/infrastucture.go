package infrastructure

import (
	"context"
	"errors"
	"fizzbuzz-code-challenge/infrastructure/logging"
	"fizzbuzz-code-challenge/infrastructure/recovery"
	"fizzbuzz-code-challenge/infrastructure/stats"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

// BuildWrapHandlerChain creates a function that wraps a handler with:
// - request counter (for stats);
// - basic request info logging;
// - basic recovery mechanism.
func BuildWrapHandlerChain(ch chan<- string) func(http.Handler) http.Handler {
	handler := stats.BuildWrapStats(ch)

	return func(next http.Handler) http.Handler {
		next = recovery.WrapRecovery(next)
		next = handler(next)
		return logging.WrapLogging(next)
	}
}

// Run runs an http server and ensures that it is gracefully shutdown:
// - in flight requests are answered;
// - new requests are not accepted.
func Run(ctx context.Context, stop func(), running func(), port int, handler http.Handler) error {
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: time.Second,
		Handler:           handler,

		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			slog.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}

		running()

		if err := httpServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}

		slog.Info("Stopped serving new connections.")
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer stopOngoingGracefully()

	return httpServer.Shutdown(shutdownCtx)
}
