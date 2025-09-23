package scanner

import (
	"fmt"
	"log/slog"

	"github.com/BookManagementSystem/scanner/mac/bluetooth"
	"github.com/BookManagementSystem/scanner/mac/common"
	"github.com/BookManagementSystem/scanner/mac/config"
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
