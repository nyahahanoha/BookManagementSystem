package store

import (
	"fmt"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
)

type BookStore interface {
	Put(book bookscommon.Info) error
	Get(isbn string) (bookscommon.Info, error)
	Del(isbn string) error
}

func NewBooksStore(config storeconfig.Config) (BookStore, error) {
	switch config.Kind {
	case storeconfig.MySQL:
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to create store: Invalid component")
	}
}
