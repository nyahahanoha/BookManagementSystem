package books

import (
	"fmt"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	googlebooks "github.com/BookManagementSystem/pkg/books/google"
	"github.com/BookManagementSystem/pkg/config"
)

type Books interface {
	Close() error
	GetInfo(isbn string) (*bookscommon.Info, error)
}

func NewBooks(config config.Config) (Books, error) {
	books, err := googlebooks.NewGoogleBooks(config.GoogleBooks)
	if err != nil {
		return nil, fmt.Errorf("failed to create books: %w", err)
	}
	return books, nil
}
