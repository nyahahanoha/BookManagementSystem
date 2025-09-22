package scanner

import (
	"fmt"
	"log/slog"

	"github.com/BookManagementSystem/pkg/scanner/bluetooth"
	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	scannerconfig "github.com/BookManagementSystem/pkg/scanner/config"
)

type Scanner interface {
	Connect() error
	Run(ch chan scannercommon.Result) error
	Close() error
}

func NewScanner(lg *slog.Logger, config scannerconfig.Config) (Scanner, error) {
	switch config.Kind {
	case scannerconfig.Bluetooth:
		bluetooth, err := bluetooth.NewBluetooth(lg, config.Bluetooth)
		if err != nil {
			return nil, fmt.Errorf("failed to create bluetooth: %w", err)
		}
		return bluetooth, nil
	case scannerconfig.Default:
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to create scanner: invalid kind")
	}
}
