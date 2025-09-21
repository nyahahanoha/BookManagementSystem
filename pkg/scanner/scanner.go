package scanner

import (
	"fmt"

	"github.com/BookManagementSystem/pkg/scanner/bluetooth"
	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	scannerconfig "github.com/BookManagementSystem/pkg/scanner/config"
)

type Scanner interface {
	Wait(ch chan scannercommon.Result) error
	Close() error
}

func NewScanner(config scannerconfig.Config) (Scanner, error) {
	switch config.Kind {
	case scannerconfig.Bluetooth:
		bluetooth, err := bluetooth.NewBluetooth(config.Bluetooth)
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
