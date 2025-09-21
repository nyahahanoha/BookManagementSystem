package config

import (
	booksconfig "github.com/BookManagementSystem/pkg/books/config"
	scannerconfig "github.com/BookManagementSystem/pkg/scanner/config"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
)

type Config struct {
	BooksConfig   booksconfig.Config   `yaml:"books"`
	StoreConfig   storeconfig.Config   `yaml:"store"`
	ScannerConfig scannerconfig.Config `yaml:"scanner"`
}
