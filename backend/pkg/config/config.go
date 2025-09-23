package config

import (
	booksconfig "github.com/BookManagementSystem/pkg/books/config"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
)

type Config struct {
	BooksConfig booksconfig.Config `yaml:"books"`
	StoreConfig storeconfig.Config `yaml:"store"`

	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}
