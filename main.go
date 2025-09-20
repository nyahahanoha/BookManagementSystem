package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	svc, err := books.NewService(ctx, option.WithAPIKey(os.Getenv("TOKEN")))
	if err != nil {
		fmt.Printf("failed to create service: %v", err)
		os.Exit(1)
	}

	volumes, err := svc.Volumes.List("isbn:9784091932518").Do()
	if err != nil {
		fmt.Printf("failed to request service: %v", err)
		os.Exit(1)
	}

	if volumes.TotalItems > 0 {
		fmt.Printf("%+v\n", volumes.Items[0])
	}
}
