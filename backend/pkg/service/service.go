package service

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	book_management_systemv1 "github.com/nyahahanoha/BookManagementSystem/backend/api/book_management_system/v1"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/books"
	bookscommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/common"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/config"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/store"
)

type BooksService struct {
	lg *slog.Logger

	books []books.Books
	store store.BookStore
}

func NewBooksService(lg *slog.Logger, config config.Config) (*BooksService, error) {
	books, err := books.NewBooks(config.BooksConfig)
	if err != nil {
		return nil, fmt.Errorf("faild to create books: %w", err)
	}

	store, err := store.NewBooksStore(lg, config.StoreConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &BooksService{
		lg:    lg.With(slog.String("Package", "service")),
		books: books,
		store: *store,
	}, nil
}

func (s *BooksService) Close() error {
	for _, b := range s.books {
		if err := b.Close(); err != nil {
			s.lg.Warn("failed to close books", slog.String("err", err.Error()))
		}
	}
	if err := s.store.Close(); err != nil {
		s.lg.Warn("failed to close store", slog.String("err", err.Error()))
		return fmt.Errorf("failed to close store: %w", err)
	}
	return nil
}

func NewAuthorizationInterceptor(token string, lg *slog.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			authHeader := req.Header().Get("Authorization")
			if authHeader != token {
				lg.Warn("unauthorized access attempt", slog.String("authHeader", authHeader))
				return nil, fmt.Errorf("unauthorized")
			}
			return next(ctx, req)
		}
	}

	return connect.UnaryInterceptorFunc(interceptor)
}

func (s *BooksService) PutBook(ctx context.Context, req *connect.Request[book_management_systemv1.PutBookRequest]) (*connect.Response[book_management_systemv1.PutBookResponse], error) {
	s.lg.Info("recieved request to Put book", slog.String("isbn", req.Msg.Isbn))
	var info *bookscommon.Info
	info, err := s.books[0].GetInfo(req.Msg.Isbn)
	if err != nil {
		s.lg.Error("failed to get info", slog.String("err", err.Error()))
	}

	for _, b := range s.books[1:] {
		moreInfo, err := b.GetInfo(req.Msg.Isbn)
		if err != nil {
			s.lg.Warn("failed to get info from another source", slog.String("err", err.Error()))
			continue
		}
		if info == nil {
			info = moreInfo
		}
		if info.Description == bookscommon.NoDescription && moreInfo.Description != bookscommon.NoDescription || info.Description == "" {
			info.Description = moreInfo.Description
		}
		if info.Image.Source.String() == "" && moreInfo.Image.Source.String() != "" {
			info.Image = moreInfo.Image
		}
		if info.Language == bookscommon.UNKOWN && moreInfo.Language != bookscommon.UNKOWN {
			info.Language = moreInfo.Language
		}
		if info.Publishdate.IsZero() && !moreInfo.Publishdate.IsZero() {
			info.Publishdate = moreInfo.Publishdate
		}
		if len(info.Authors) == 0 && len(moreInfo.Authors) != 0 {
			info.Authors = moreInfo.Authors
		}
		if info.Title == "" && moreInfo.Title != "" {
			info.Title = moreInfo.Title
		}
	}
	if info.Title == "" {
		info.Title = req.Msg.Isbn
	}
	s.lg.Info("Get book info", slog.String("title", info.Title), slog.String("isbn", info.ISBN))
	if err := s.store.Put(*info); err != nil {
		s.lg.Error("failed to put info in store", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to put info in store: %w", err)
	}
	return connect.NewResponse(&book_management_systemv1.PutBookResponse{
		Book: convertInfoToProtobuf(*info),
	}), nil
}

func convertInfoToProtobuf(info bookscommon.Info) *book_management_systemv1.Book {
	var language book_management_systemv1.Language
	switch info.Language {
	case bookscommon.JP:
		language = book_management_systemv1.Language_JAPANESE
	case bookscommon.EN:
		language = book_management_systemv1.Language_ENGLISH
	default:
		language = book_management_systemv1.Language_UNKNOWN
	}
	return &book_management_systemv1.Book{
		Isbn:        info.ISBN,
		Title:       info.Title,
		Authors:     info.Authors,
		Description: info.Description,
		Publishdate: info.Publishdate.Format("2006-01"),
		Language:    language,
		Imageurl:    info.Image.Path,
	}
}

func (s *BooksService) GetBook(ctx context.Context, req *connect.Request[book_management_systemv1.GetBookRequest]) (*connect.Response[book_management_systemv1.GetBookResponse], error) {
	s.lg.Info("recieved request to Get book", slog.String("isbn", req.Msg.Isbn))
	info, err := s.store.Get(req.Msg.Isbn)
	if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to get book in store: %w", err)
	}

	return connect.NewResponse(&book_management_systemv1.GetBookResponse{
		Book: convertInfoToProtobuf(info),
	}), nil
}

func (s *BooksService) GetAllBooks(ctx context.Context, req *connect.Request[book_management_systemv1.GetAllBooksRequest]) (*connect.Response[book_management_systemv1.GetAllBooksResponse], error) {
	s.lg.Info("recieved request to Get all books")
	books, err := s.store.GetAll()
	if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to get all books in store: %w", err)
	}

	return connect.NewResponse(&book_management_systemv1.GetAllBooksResponse{
		Books: func() []*book_management_systemv1.Book {
			res := make([]*book_management_systemv1.Book, 0, len(books))
			for _, info := range books {
				res = append(res, convertInfoToProtobuf(info))
			}
			return res
		}(),
	}), nil
}

func (s *BooksService) SearchBook(ctx context.Context, req *connect.Request[book_management_systemv1.SearchBookRequest]) (*connect.Response[book_management_systemv1.SearchBookResponse], error) {
	s.lg.Info("recieved request to Search books", slog.String("title", req.Msg.Title))
	books, err := s.store.Search(req.Msg.Title)
	if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to search books in store: %w", err)
	}
	return connect.NewResponse(&book_management_systemv1.SearchBookResponse{
		Books: func() []*book_management_systemv1.Book {
			res := make([]*book_management_systemv1.Book, 0, len(books))
			for _, info := range books {
				res = append(res, convertInfoToProtobuf(info))
			}
			return res
		}(),
	}), nil
}

func (s *BooksService) DeleteBook(ctx context.Context, req *connect.Request[book_management_systemv1.DeleteBookRequest]) (*connect.Response[book_management_systemv1.DeleteBookResponse], error) {
	s.lg.Info("recieved request to Delete book", slog.String("isbn", req.Msg.Isbn))
	if err := s.store.Del(req.Msg.Isbn); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to delete book in store: %w", err)
	}
	return connect.NewResponse(&book_management_systemv1.DeleteBookResponse{}), nil
}

func (s *BooksService) RenameBook(ctx context.Context, req *connect.Request[book_management_systemv1.RenameBookRequest]) (*connect.Response[book_management_systemv1.RenameBookResponse], error) {
	s.lg.Info("recieved request to Rename book", slog.String("isbn", req.Msg.Isbn), slog.String("title", req.Msg.Title))
	if err := s.store.Rename(req.Msg.Isbn, req.Msg.Title); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to rename book in store: %w", err)
	}
	return connect.NewResponse(&book_management_systemv1.RenameBookResponse{}), nil
}
