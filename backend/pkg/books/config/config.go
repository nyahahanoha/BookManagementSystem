package booksconfig

type Config struct {
	Kind   []BooksComponent  `yaml:"kind"`
	Google GoogleBooksConfig `yaml:"google"`
}

//go:generate go run github.com/dmarkham/enumer -type=BooksComponent -yaml
type BooksComponent uint32

const (
	Google BooksComponent = iota
	NDL
)

type GoogleBooksConfig struct {
	APIKey string `yaml:"api_key"`
}
