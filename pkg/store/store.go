package store

import (
	"fmt"
	"log/slog"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
	"github.com/BookManagementSystem/pkg/store/db"
	"github.com/BookManagementSystem/pkg/store/object"
)

type BookStore struct {
	db     db.DBStore
	object object.ObjectStore
}

func NewBooksStore(lg *slog.Logger, config storeconfig.Config) (*BookStore, error) {
	db, err := db.NewDBStore(lg, config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	object, err := object.NewObjectStore(lg, config.Object)
	if err != nil {
		return nil, fmt.Errorf("failed to connect object: %w", err)
	}

	return &BookStore{
		db:     db,
		object: object,
	}, nil
}

func (s *BookStore) Put(book bookscommon.Info) error {
	if err := s.db.Put(book); err != nil {
		return fmt.Errorf("failed to put info in db: %w", err)
	}
	if err := s.object.Put(book.Image.Source, book.ISBN); err != nil {
		return fmt.Errorf("failed to put image in object: %w", err)
	}
	return nil
}
func (s *BookStore) Get(isbn string) (*bookscommon.Info, error) {
	info, err := s.db.Get(isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get info in db: %w", err)
	}
	path, err := s.object.Get(isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get image in object: %w", err)
	}
	info.Image.Path = path
	return info, err
}
func (s *BookStore) Del(isbn string) error { return nil }
