package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"proektUrl/internal/httpapi"
	"proektUrl/internal/link"
	"syscall"
	"time"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	addr := ":8080"
	host := "http://localhost" + addr

	storage := link.NewMemoryStorage()
	svc := link.NewService(storage)
	handler := httpapi.NewHandler(svc, host, logger)
	router := httpapi.NewRouter(handler)

	srv := &http.Server{

		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {

		logger.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			logger.Error("server falied", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown falied", "err", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}
