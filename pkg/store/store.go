package store

import (
	"fmt"
	"log/slog"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
	"github.com/BookManagementSystem/pkg/store/db"
)

type BookStore struct {
	db db.DBStore
}

func NewBooksStore(lg *slog.Logger, config storeconfig.Config) (*BookStore, error) {
	db, err := db.NewDBStore(lg, config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}
	return &BookStore{
		db: db,
	}, nil
}

func (s *BookStore) Put(book bookscommon.Info) error {
	if err := s.db.Put(book); err != nil {
		return fmt.Errorf("failed to put info in db: %w", err)
	}
	return nil
}
func (s *BookStore) Get(isbn string) (*bookscommon.Info, error) {
	info, err := s.db.Get(isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get info in db: %w", err)
	}
	return info, nil
}
func (s *BookStore) Del(isbn string) error { return nil }
