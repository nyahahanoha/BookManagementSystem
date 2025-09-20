package books

import (
	"fmt"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	booksconfig "github.com/BookManagementSystem/pkg/books/config"
	googlebooks "github.com/BookManagementSystem/pkg/books/google"
)

type Books interface {
	Close() error
	GetInfo(isbn string) (*bookscommon.Info, error)
}

func NewBooks(config booksconfig.Config) (Books, error) {
	switch config.Kind {
	case booksconfig.Google:
		books, err := googlebooks.NewGoogleBooks(config.Google)
		if err != nil {
			return nil, fmt.Errorf("failed to create books: %w", err)
		}
		return books, nil
	default:
		return nil, fmt.Errorf("failed to create books: Unknown Compoent")
	}
}
