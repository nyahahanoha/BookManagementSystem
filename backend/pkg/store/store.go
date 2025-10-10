package store

import (
	"fmt"
	"log/slog"

	bookscommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/common"
	storecommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/common"
	storeconfig "github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/config"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/db"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/object"
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
func (s *BookStore) Get(isbn string) (bookscommon.Info, error) {
	info, err := s.db.Get(isbn)
	if err == storecommon.ErrNotFoundBook {
		return bookscommon.Info{}, err
	} else if err != nil {
		return bookscommon.Info{}, fmt.Errorf("failed to get info in db: %w", err)
	}
	path, err := s.object.Get(isbn)
	if err != nil {
		return bookscommon.Info{}, fmt.Errorf("failed to get image in object: %w", err)
	}
	info.Image.Path = path
	return info, nil
}

func (s *BookStore) GetAll() ([]bookscommon.Info, error) {
	books, err := s.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get info in db: %w", err)
	}
	for i, book := range books {
		path, err := s.object.Get(book.ISBN)
		if err != nil {
			return nil, fmt.Errorf("failed to get image in object: %w", err)
		}
		books[i].Image.Path = path
	}
	return books, nil
}

func (s *BookStore) Search(title string) ([]bookscommon.Info, error) {
	books, err := s.db.Search(title)
	if err != nil {
		return nil, fmt.Errorf("failed to get info in db: %w", err)
	}
	for i, book := range books {
		path, err := s.object.Get(book.ISBN)
		if err != nil {
			return nil, fmt.Errorf("failed to get image in object: %w", err)
		}
		books[i].Image.Path = path
	}
	return books, nil
}

func (s *BookStore) Del(isbn string) error {
	if err := s.db.Delete(isbn); err != nil {
		return fmt.Errorf("failed to delete info in db: %w", err)
	}
	if err := s.object.Delete(isbn); err != nil {
		return fmt.Errorf("failed to delete image in object: %w", err)
	}
	return nil
}
