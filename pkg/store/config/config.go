package storeconfig

type Config struct {
	DB     DBConfig     `yaml:"db"`
	Object ObjectConfig `yaml:"object"`
}
