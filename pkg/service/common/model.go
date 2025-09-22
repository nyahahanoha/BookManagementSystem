package servicecommon

import bookscommon "github.com/BookManagementSystem/pkg/books/common"

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type BooksResponse struct {
	Books []bookscommon.Info `json:"books"`
	Count int                `json:"count"`
}
