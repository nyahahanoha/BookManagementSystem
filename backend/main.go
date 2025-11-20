package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"connectrpc.com/connect"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/nyahahanoha/BookManagementSystem/backend/api/book_management_system/v1/book_management_systemv1connect"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/config"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/service"
	"github.com/rs/cors"
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

	jwks, err := jwk.Fetch(context.Background(), cfg.PomeriumJWKSURL)
	if err != nil {
		logger.Error("failed to fetch jwks", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pomeriumAuth := service.NewAuthInterceptor(jwks, logger, cfg.AdminEmail)

	mux := http.NewServeMux()
	path, handler := book_management_systemv1connect.NewBookManagementServiceHandler(
		books,
		connect.WithInterceptors(pomeriumAuth),
	)
	mux.Handle(path, handler)

	mux.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/images/") {
			filePath := "/var/lib/booksystem" + r.URL.Path[len("/images"):]
			http.ServeFile(w, r, filePath)
		}
	})

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{cfg.FrontendURL}, // フロントのURL
		AllowedMethods:     []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Content-Type", "X-Pomerium-Jwt-Assertion"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
	})

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
		Handler: c.Handler(mux),
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
