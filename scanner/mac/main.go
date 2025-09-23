package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BookManagementSystem/scanner/mac/bluetooth"
	"github.com/BookManagementSystem/scanner/mac/common"
	"github.com/BookManagementSystem/scanner/mac/config"
	"go.yaml.in/yaml/v2"
)

type Scanner interface {
	Connect() error
	Run(ch chan common.Result) error
	Close() error
}

func NewScanner(lg *slog.Logger, c config.Config) (Scanner, error) {
	switch c.Kind {
	case config.Bluetooth:
		bluetooth, err := bluetooth.NewBluetooth(lg, c.Bluetooth)
		if err != nil {
			return nil, fmt.Errorf("failed to create bluetooth: %w", err)
		}
		return bluetooth, nil
	case config.Default:
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to create scanner: invalid kind")
	}
}

func main() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to reaf yaml file: %v", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("failed to unmarshal yaml file: %v", err)
	}

	if cfg.API == "" || cfg.Token == "" {
		log.Fatalf("api or token is empty")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	scanner, err := NewScanner(logger, cfg)
	if err != nil {
		log.Fatalf("failed to create scanner: %v", err)
	}

	if err := scanner.Connect(); err != nil {
		log.Fatalf("failed to connect scanner: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ch := make(chan common.Result, 1)

	go func() {
		if err := scanner.Run(ch); err != nil {
			log.Fatal("failed to run scanner: %w", err)
		}
	}()

	for {
		select {
		case result := <-ch:
			req, err := http.NewRequest("POST", cfg.API+"/put:"+result.ISBN, nil)
			if err != nil {
				logger.Error("failed to post request", slog.String("isbn", result.ISBN), slog.String("error", err.Error()))
				continue
			}
			req.Header.Set("Authorization", cfg.Token)
			client := &http.Client{}
			_, err = client.Do(req)
			if err != nil {
				logger.Error("failed to post request", slog.String("isbn", result.ISBN), slog.String("error", err.Error()))
				continue
			}
		case <-sigs:
			if err := scanner.Close(); err != nil {
				log.Fatalf("failed to close scanner: %v", err)
				return
			}
			return
		}
	}
}
