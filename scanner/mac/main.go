package scanner

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/BookManagementSystem/scanner/mac/scanner/bluetooth"
	"github.com/BookManagementSystem/scanner/mac/scanner/common"
	"github.com/BookManagementSystem/scanner/mac/scanner/config"
	"go.yaml.in/yaml/v2"
)

type Scanner interface {
	Connect() error
	Run(ch chan common.Result) error
	Close() error
}

func NewScanner(lg *slog.Logger, config config.Config) (Scanner, error) {
	switch config.Kind {
	case config.Bluetooth:
		bluetooth, err := bluetooth.NewBluetooth(lg, config.Bluetooth)
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	scanner, err := NewScanner(logger, cfg)
	if err != nil {
		log.Fatalf("failed to create scanner: %v", err)
	}

	if err := scanner.Connect(); err != nil {
		log.Fatalf("failed to connect scanner: %v", err)
	}
	defer func() {
		if err := scanner.Close(); err != nil {
			log.Fatalf("failed to close scanner: %v", err)
		}
	}()
}
