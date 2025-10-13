package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	book_management_systemv1 "github.com/nyahahanoha/BookManagementSystem/api/book_management_system/v1"
	book_management_systemv1connect "github.com/nyahahanoha/BookManagementSystem/api/book_management_system/v1/book_management_systemv1connect"

	"connectrpc.com/connect"
	"github.com/BookManagementSystem/scanner/mac/authorization"
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

	file = []byte(os.ExpandEnv(string(file)))

	var cfg config.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("failed to unmarshal yaml file: %v", err)
	}

	if cfg.API == "" || cfg.CallBackPort == 0 {
		log.Fatalf("api or callback port is empty")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	apiURI, err := url.Parse(cfg.API)
	if err != nil {
		log.Fatalf("failed to parse api url: %v", err)
	}
	token, err := authorization.Authorization(*apiURI, cfg.CallBackPort)
	if err != nil {
		log.Fatal("failed to get token: %w", err)
	}

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

	client := book_management_systemv1connect.NewBookManagementServiceClient(
		http.DefaultClient,
		cfg.API,
	)
	ctx, cancel := context.WithCancel(context.Background())
	ctx, callInfo := connect.NewClientContext(ctx)
	callInfo.RequestHeader().Set("Authorization", "Pomerium "+token)
	for {
		select {
		case result := <-ch:
			if _, err = client.PutBook(ctx, connect.NewRequest(&book_management_systemv1.PutBookRequest{
				Isbn: result.ISBN,
			})); err != nil {
				logger.Error("failed to post request", slog.String("isbn", result.ISBN), slog.String("error", err.Error()))
				continue
			}
		case <-sigs:
			if err := scanner.Close(); err != nil {
				log.Fatalf("failed to close scanner: %v", err)
				return
			}
			cancel()
			return
		}
	}

}
