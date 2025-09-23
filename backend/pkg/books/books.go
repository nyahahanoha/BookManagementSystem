package books

import (
	"fmt"

	bookscommon "github.com/BookManagementSystem/backend/pkg/books/common"
	booksconfig "github.com/BookManagementSystem/backend/pkg/books/config"
	googlebooks "github.com/BookManagementSystem/backend/pkg/books/google"
	ndlbooks "github.com/BookManagementSystem/backend/pkg/books/ndl"
)

type Books interface {
	Close() error
	GetInfo(isbn string) (*bookscommon.Info, error)
}

func NewBooks(config booksconfig.Config) ([]Books, error) {
	var booksList []Books
	for _, kind := range config.Kind {
		switch kind {
		case booksconfig.Google:
			books, err := googlebooks.NewGoogleBooks(config.Google)
			if err != nil {
				return nil, fmt.Errorf("failed to create books: %w", err)
			}
			booksList = append(booksList, books)
		case booksconfig.NDL:
			books, err := ndlbooks.NewNDL()
			if err != nil {
				return nil, fmt.Errorf("failed to create books: %w", err)
			}
			booksList = append(booksList, books)
		default:
			return nil, fmt.Errorf("failed to create books: Unknown Compoent")
		}
	}
	return booksList, nil
}
