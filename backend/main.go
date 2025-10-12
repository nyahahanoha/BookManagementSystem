package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	"github.com/nyahahanoha/BookManagementSystem/api/book_management_system/v1/book_management_systemv1connect"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/config"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/service"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to reaf yaml file: %v", err)
	}

	file = []byte(os.ExpandEnv(string(file)))

	var cfg config.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("failed to unmarshal yaml file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	books, err := service.NewBooksService(logger, cfg)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}
	mux := http.NewServeMux()
	path, handler := book_management_systemv1connect.NewBookManagementServiceHandler(
		books,
		connect.WithInterceptors(service.NewAuthorizationInterceptor(cfg.Token, logger)),
	)
	mux.Handle(path, handler)

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Info("received signal, shutting down", slog.String("signal", sig.String()))
		cancel()
	}()

	logger.Info("starting server", slog.String("path", path))
	defer func() {
		if err := books.Close(); err != nil {
			logger.Error("failed to close service", slog.String("err", err.Error()))
		}
		logger.Info("server stopped")
	}()

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("failed to shutdown server", slog.String("err", err.Error()))
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", slog.String("err", err.Error()))
	}
}
