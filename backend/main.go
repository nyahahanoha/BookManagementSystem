package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/BookManagementSystem/backend/pkg/config"
	"github.com/BookManagementSystem/backend/pkg/service"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to reaf yaml file: %v", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("failed to unmarshal yaml file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	service, err := service.NewBooksService(logger, cfg)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	go func() {
		if err := service.Listen(); err != nil {
			logger.Error("failed to service", slog.String("err", err.Error()))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	if err := service.Shutdown(); err != nil {
		log.Fatalf("failed to close service: %v", err)
	}
}
