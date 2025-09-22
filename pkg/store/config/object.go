package storeconfig

type ObjectConfig struct {
	Kind       ObjectComponent `yaml:"kind"`
	FileConfig FileConfig      `yaml:"file"`
}

//go:generate go run github.com/dmarkham/enumer -type=ObjectComponent -yaml
type ObjectComponent uint32

const (
	FileSystem ObjectComponent = iota
)

type FileConfig struct {
	Prefix string `yaml:"path"`
}
