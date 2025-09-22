package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/BookManagementSystem/pkg/books"
	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	"github.com/BookManagementSystem/pkg/config"
	"github.com/BookManagementSystem/pkg/scanner"
	scannercommon "github.com/BookManagementSystem/pkg/scanner/common"
	servicecommon "github.com/BookManagementSystem/pkg/service/common"
	"github.com/BookManagementSystem/pkg/store"
	storecommon "github.com/BookManagementSystem/pkg/store/common"
)

type BooksService struct {
	lg *slog.Logger

	books   books.Books
	scanner scanner.Scanner
	store   store.BookStore
	api     *http.Server
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
		lg:      lg.With(slog.String("Package", "service")),
		books:   books,
		scanner: scanner,
		store:   *store,
		api: &http.Server{
			Addr: config.Address,
		},
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
					return
				}
				if err := s.store.Put(*info); err != nil {
					s.lg.Error("failed to put info in store", slog.String("err", err.Error()))
				}
				book, err := s.store.Get(result.ISBN)
				if err != nil {
					s.lg.Error("failed to get info in store", slog.String("err", err.Error()))
				}
				s.lg.Info("Get Book Info",
					slog.Any("info", book),
				)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.api.Shutdown(ctx); err != nil {
		s.lg.Info("Server closed")
	}
	return nil
}

func CORSMiddleware() rest.Middleware {
	return rest.MiddlewareSimple(func(handler rest.HandlerFunc) rest.HandlerFunc {
		return func(w rest.ResponseWriter, r *rest.Request) {
			// CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler(w, r)
		}
	})
}

func (s *BooksService) Listen() error {
	api := rest.NewApi()
	api.Use(CORSMiddleware())
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/books", s.GetAllBooks),
		rest.Get("/book:isbn", s.GetBook),
		rest.Get("/books/search:title", s.SearchBook),
	)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}
	api.SetApp(router)
	s.api.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/images/") {
			filePath := "/var/lib/booksystem" + r.URL.Path[len("/images"):]
			http.ServeFile(w, r, filePath)
			return
		}
		api.MakeHandler().ServeHTTP(w, r)
	})
	if err := s.api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}

	s.lg.Info("Close Listen Service")
	return nil
}

func (s *BooksService) GetBook(w rest.ResponseWriter, r *rest.Request) {
	isbn := strings.ReplaceAll(r.PathParam("isbn"), ":", "")

	info, err := s.store.Get(isbn)
	if err == storecommon.ErrNotFoundBook {
		if err := w.WriteJson(servicecommon.BooksResponse{
			Books: nil,
			Count: 0,
		}); err != nil {
			s.lg.Error("internal server error", slog.String("err", err.Error()))
			rest.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	} else if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := w.WriteJson(servicecommon.BooksResponse{
		Books: []bookscommon.Info{info},
		Count: 1,
	}); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *BooksService) GetAllBooks(w rest.ResponseWriter, r *rest.Request) {
	books, err := s.store.GetAll()
	if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := w.WriteJson(servicecommon.BooksResponse{
		Books: books,
		Count: len(books),
	}); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *BooksService) SearchBook(w rest.ResponseWriter, r *rest.Request) {
	rawTitle := r.PathParam("title")
	decodedTitle, err := url.PathUnescape(rawTitle)
	if err != nil {
		decodedTitle = rawTitle
	}
	title := strings.ReplaceAll(decodedTitle, ":", "")

	books, err := s.store.Search(title)
	if err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := w.WriteJson(servicecommon.BooksResponse{
		Books: books,
		Count: len(books),
	}); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
