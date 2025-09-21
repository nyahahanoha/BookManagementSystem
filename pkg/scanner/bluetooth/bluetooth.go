package bluetooth

import (
	"fmt"
	"time"

	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	scannerconfig "github.com/BookManagementSystem/pkg/scanner/config"
	ble "tinygo.org/x/bluetooth"
)

type Bluetooth struct {
	ServiceUUID        string
	CharacteristicUUID string

	device  ble.Device
	closech chan struct{}
}

var adapter = ble.DefaultAdapter

func NewBluetooth(config scannerconfig.BluetoothConfig) (*Bluetooth, error) {
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("failed to create bluetooth adapater: %w", err)
	}
	ch := make(chan ble.ScanResult, 1)
	if err := adapter.Scan(func(adapter *ble.Adapter, result ble.ScanResult) {
		if result.LocalName() == config.Name {
			adapter.StopScan()
			ch <- result
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to scan bluetooth scanner: %w", err)
	}

	select {
	case result := <-ch:
		device, err := adapter.Connect(result.Address, ble.ConnectionParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect bluetooth scanner: %w", err)
		}

		return &Bluetooth{
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
	close(s.closech)
	s.device.Disconnect()
	return nil
}

func (s *Bluetooth) Wait(ch chan scannercommon.Result) error {
	services, err := s.device.DiscoverServices(nil)
	if err != nil {
		return fmt.Errorf("failed to discover service: %w", err)
	}

	for _, service := range services {
		if service.UUID().String() == s.ServiceUUID {
			chars, err := service.DiscoverCharacteristics(nil)
			if err != nil {
				return fmt.Errorf("failed to discover characteristics: %w", err)
			}

			for _, char := range chars {
				if char.UUID().String() == s.CharacteristicUUID {
					if err = char.EnableNotifications(func(buf []byte) {
						ch <- scannercommon.Result{
							ISBN: string(buf),
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
