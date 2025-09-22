package storeconfig

type DBConfig struct {
	Kind        DBComponent `yaml:"kind"`
	MySQLConfig MySQLConfig `yaml:"mysql"`
}

//go:generate go run github.com/dmarkham/enumer -type=DBComponent -yaml
type DBComponent uint32

const (
	MySQL DBComponent = iota
	PostgreSQL
)

type MySQLConfig struct {
	IPAddress string `yaml:"address"`
	Port      uint16 `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
}
