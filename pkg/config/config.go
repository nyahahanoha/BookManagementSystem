package config

import booksconfig "github.com/BookManagementSystem/pkg/books/config"

type Config struct {
	BooksConfig booksconfig.Config `yaml:"books"`
}
