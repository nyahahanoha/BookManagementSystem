package storeconfig

type ObjectConfig struct {
	Kind       ObjectComponent `yaml:"kind"`
	FileConfig FileConfig      `yaml:"file"`
}

type ObjectComponent uint32

const (
	FileSystem ObjectComponent = iota
)

type FileConfig struct {
	Prefix string `yaml:"path"`
}
