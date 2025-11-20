package scanner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nyahahanoha/BookManagementSystem/scanner/mac/bluetooth"
	"github.com/nyahahanoha/BookManagementSystem/scanner/mac/common"
	"github.com/nyahahanoha/BookManagementSystem/scanner/mac/config"
)

type Scanner interface {
	Connect() error
	Run(ctx context.Context, ch chan common.Result) error
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
