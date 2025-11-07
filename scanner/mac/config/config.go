package config

import "time"

type Config struct {
	Kind         ScannerComponent `yaml:"kind"`
	API          string           `yaml:"api"`
	CallBackPort uint16           `yaml:"callback_port"`
	Default      DefaultConfig    `yaml:"default"`
	Bluetooth    BluetoothConfig  `yaml:"bluetooth"`
}

//go:generate go run github.com/dmarkham/enumer -type=ScannerComponent -yaml
type ScannerComponent uint32

const (
	Default ScannerComponent = iota
	Bluetooth
)

type DefaultConfig struct {
}

type BluetoothConfig struct {
	Enabled            bool          `yaml:"enabled"`
	Name               string        `yaml:"name"`
	Timeout            time.Duration `yaml:"timeout"`
	ServiceUUID        string        `yaml:"service"`
	CharacteristicUUID string        `yaml:"characteristic"`
}
