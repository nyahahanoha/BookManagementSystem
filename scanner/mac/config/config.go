package config

import "time"

type Config struct {
	Kind      ScannerComponent `yaml:"kind"`
	API       string           `yaml:"api"`
	Token     string           `yaml:"token"`
	Default   DefaultConfig    `yaml:"default"`
	Bluetooth BluetoothConfig  `yaml:"bluetooth"`
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
	Name               string        `yaml:"name"`
	Timeout            time.Duration `yaml:"timeout"`
	ServiceUUID        string        `yaml:"service"`
	CharacteristicUUID string        `yaml:"characteristic"`
}
