package bluetooth

import (
	"fmt"
	"log/slog"
	"time"

	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	scannerconfig "github.com/BookManagementSystem/pkg/scanner/config"
	"github.com/dlclark/regexp2"
	ble "tinygo.org/x/bluetooth"
)

type Bluetooth struct {
	lg *slog.Logger

	ServiceUUID        string
	CharacteristicUUID string

	device  ble.Device
	closech chan struct{}
}

var adapter = ble.DefaultAdapter

func NewBluetooth(lg *slog.Logger, config scannerconfig.BluetoothConfig) (*Bluetooth, error) {
	lg.With(slog.String("Package", "bluetooth"))
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("failed to create bluetooth adapater: %w", err)
	}

	ch := make(chan ble.ScanResult, 1)
	if err := adapter.Scan(func(adapter *ble.Adapter, result ble.ScanResult) {
		lg.Info("found device",
			slog.String("name", result.LocalName()),
			slog.String("address", result.Address.String()),
		)
		if result.LocalName() == config.Name {
			if err := adapter.StopScan(); err != nil {
				lg.Error("failed to stop scan")
			}
			ch <- result
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to scan bluetooth scanner: %w", err)
	}

	select {
	case result := <-ch:
		lg.Info("Connecting...",
			slog.String("device", result.LocalName()),
		)
		device, err := adapter.Connect(result.Address, ble.ConnectionParams{})
		lg.Info("Connected!",
			slog.String("device", result.LocalName()),
		)

		if err != nil {
			return nil, fmt.Errorf("failed to connect bluetooth scanner: %w", err)
		}

		return &Bluetooth{
			lg:                 lg,
			ServiceUUID:        config.ServiceUUID,
			CharacteristicUUID: config.CharacteristicUUID,
			device:             device,
			closech:            make(chan struct{}),
		}, nil

	case <-time.After(config.Timeout):
		return nil, fmt.Errorf("failed to connect bluetooth scanner: Timeout")
	}
}

func (s *Bluetooth) Close() error {
	s.lg.Info("Close Bluetooth")
	close(s.closech)
	return s.device.Disconnect()
}

func (s *Bluetooth) Wait(ch chan scannercommon.Result) error {
	services, err := s.device.DiscoverServices(nil)
	if err != nil {
		return fmt.Errorf("failed to discover service: %w", err)
	}

	re := regexp2.MustCompile(`\d+`, 0)

	for _, service := range services {
		if service.UUID().String() == s.ServiceUUID {
			chars, err := service.DiscoverCharacteristics(nil)
			if err != nil {
				return fmt.Errorf("failed to discover characteristics: %w", err)
			}

			for _, char := range chars {
				if char.UUID().String() == s.CharacteristicUUID {
					if err = char.EnableNotifications(func(buf []byte) {
						isbn, err := re.FindStringMatch(string(buf))
						if err != nil {
							s.lg.Error("failed to regexp", slog.String("err", err.Error()))
						}
						s.lg.Info("get ISBN", slog.String("isbn", isbn.String()))
						ch <- scannercommon.Result{
							ISBN: isbn.String(),
						}

					}); err != nil {
						return fmt.Errorf("failed to notification: %w", err)
					}
				}
			}
		}
	}

	<-s.closech
	close(ch)

	return nil
}
