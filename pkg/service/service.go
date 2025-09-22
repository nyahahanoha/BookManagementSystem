package service

import (
	"fmt"
	"log/slog"

	"github.com/BookManagementSystem/pkg/books"
	"github.com/BookManagementSystem/pkg/config"
	"github.com/BookManagementSystem/pkg/scanner"
	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	"github.com/BookManagementSystem/pkg/store"
)

type BooksService struct {
	lg      *slog.Logger
	books   books.Books
	scanner scanner.Scanner
	store   store.BookStore
}

func NewBooksService(lg *slog.Logger, config config.Config) (*BooksService, error) {
	books, err := books.NewBooks(config.BooksConfig)
	if err != nil {
		return nil, fmt.Errorf("faild to create books: %w", err)
	}

	scanner, err := scanner.NewScanner(lg, config.ScannerConfig)
	if err != nil {
		return nil, fmt.Errorf("faild to create scanner: %w", err)
	}

	store, err := store.NewBooksStore(lg, config.StoreConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &BooksService{
		books:   books,
		scanner: scanner,
		store:   *store,
		lg:      lg.With(slog.String("Package", "service")),
	}, nil
}

func (s *BooksService) Run() error {
	s.lg.Info("Running BooksService")
	ch := make(chan scannercommon.Result)

	go func() {
		for {
			result, ok := <-ch
			if !ok {
				s.lg.Info("Stop running BooksService")
				return
			}
			go func(result scannercommon.Result) {
				info, err := s.books.GetInfo(result.ISBN)
				if err != nil {
					s.lg.Error("failed to get info", slog.String("err", err.Error()))
				}
				s.lg.Info("Get Book Info",
					slog.Any("info", info),
				)
				if err := s.store.Put(*info); err != nil {
					s.lg.Error("failed to put info in store", slog.String("err", err.Error()))
				}
			}(result)
		}
	}()

	if err := s.scanner.Run(ch); err != nil {
		return fmt.Errorf("failed to running scanner: %w", err)
	}
	return nil
}

func (s *BooksService) Close() error {
	s.lg.Info("Close BooksService")
	if err := s.books.Close(); err != nil {
		s.lg.Error("Failed to close books", slog.String("err", err.Error()))
	}
	if err := s.scanner.Close(); err != nil {
		s.lg.Error("Failed to close scanner", slog.String("err", err.Error()))
	}
	return nil
}
