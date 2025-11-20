package main

import (
	"bufio"
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
	"github.com/BookManagementSystem/scanner/mac/common"
	"github.com/BookManagementSystem/scanner/mac/config"
	"github.com/BookManagementSystem/scanner/mac/scanner"
	"go.yaml.in/yaml/v2"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ch := make(chan common.Result, 1)

	client := book_management_systemv1connect.NewBookManagementServiceClient(
		http.DefaultClient,
		cfg.API,
	)
	ctx, cancel := context.WithCancel(context.Background())
	ctx, callInfo := connect.NewClientContext(ctx)
	callInfo.RequestHeader().Set("Authorization", "Pomerium "+token)

	go func() {
		if !cfg.Bluetooth.Enabled {
			return
		}
		scanner, err := scanner.NewScanner(logger, cfg)
		if err != nil {
			log.Fatalf("failed to create scanner: %v", err)
		}

		if err := scanner.Connect(); err != nil {
			log.Fatalf("failed to connect scanner: %v", err)
		}

		if err := scanner.Run(ctx, ch); err != nil {
			log.Fatal("failed to run scanner: %w", err)
		}
	}()

	go func() {
		stdscanner := bufio.NewScanner(os.Stdin)
		var isbn string
		for {
			fmt.Print("\nPlease input ISBN: ")
			if !stdscanner.Scan() {
				break
			}
			isbn = stdscanner.Text()
			if isbn != "" {
				ch <- common.Result{ISBN: isbn}
			}
		}
	}()

	for {
		select {
		case result := <-ch:
			if res, err := client.PutBook(ctx, connect.NewRequest(&book_management_systemv1.PutBookRequest{
				Isbn: result.ISBN,
			})); err != nil {
				logger.Error("failed to put request", slog.String("isbn", result.ISBN), slog.String("error", err.Error()))
				continue
			} else {
				logger.Info("succeeded to put book", slog.String("isbn", result.ISBN))
				log.Printf("Response Book Info: %+v", res.Msg.GetBook())
			}
		case <-sigs:
			cancel()
			return
		}
	}
}
