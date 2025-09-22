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

	Name               string
	Timeout            time.Duration
	ServiceUUID        string
	CharacteristicUUID string

	device  *ble.Device
	closech chan struct{}
}

var adapter = ble.DefaultAdapter

func NewBluetooth(lg *slog.Logger, config scannerconfig.BluetoothConfig) (*Bluetooth, error) {
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("failed to create bluetooth adapater: %w", err)
	}
	return &Bluetooth{
		lg:                 lg.With(slog.String("Package", "bluetooth")),
		Name:               config.Name,
		Timeout:            config.Timeout,
		ServiceUUID:        config.ServiceUUID,
		CharacteristicUUID: config.CharacteristicUUID,
	}, nil
}

func (s *Bluetooth) Close() error {
	s.lg.Info("Close Bluetooth")
	if s.device != nil {
		close(s.closech)
		return s.device.Disconnect()
	}
	return nil
}

func (s *Bluetooth) Connect() error {
	ch := make(chan ble.ScanResult, 1)
	if err := adapter.Scan(func(adapter *ble.Adapter, result ble.ScanResult) {
		s.lg.Info("found device",
			slog.String("name", result.LocalName()),
			slog.String("address", result.Address.String()),
		)
		if result.LocalName() == s.Name {
			if err := adapter.StopScan(); err != nil {
				s.lg.Error("failed to stop scan")
			}
			ch <- result
		}
	}); err != nil {
		return fmt.Errorf("failed to scan bluetooth scanner: %w", err)
	}

	select {
	case result := <-ch:
		s.lg.Info("Connecting...",
			slog.String("device", result.LocalName()),
		)
		device, err := adapter.Connect(result.Address, ble.ConnectionParams{})
		s.lg.Info("Connected!",
			slog.String("device", result.LocalName()),
		)

		if err != nil {
			return fmt.Errorf("failed to connect bluetooth scanner: %w", err)
		}

		s.device = &device
		return nil
	case <-time.After(s.Timeout):
		return fmt.Errorf("failed to connect bluetooth scanner: Timeout")
	}
}

func (s *Bluetooth) Run(ch chan scannercommon.Result) error {
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
	s.closech = make(chan struct{})

	<-s.closech
	close(ch)

	return nil
}
