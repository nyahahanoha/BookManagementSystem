package store

import (
	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
)

type BookStore struct {
}

func NewBooksStore(config storeconfig.Config) (BookStore, error) {
	return BookStore{}, nil
}

func (s *BookStore) Put(book bookscommon.Info) error            { return nil }
func (s *BookStore) Get(isbn string) (*bookscommon.Info, error) { return nil, nil }
func (s *BookStore) Del(isbn string) error                      { return nil }
