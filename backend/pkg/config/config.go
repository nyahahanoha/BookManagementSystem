package config

import (
	booksconfig "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/config"
	storeconfig "github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/config"
)

type Config struct {
	BooksConfig booksconfig.Config `yaml:"books"`
	StoreConfig storeconfig.Config `yaml:"store"`

	Address         string `yaml:"address"`
	AdminEmail      string `yaml:"admin_email"`
	PomeriumJWKSURL string `yaml:"pomerium_jwks_url"`
	FrontendURL     string `yaml:"frontend_url"`
}
