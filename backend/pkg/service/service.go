package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BookManagementSystem/backend/pkg/books"
	bookscommon "github.com/BookManagementSystem/backend/pkg/books/common"
	"github.com/BookManagementSystem/backend/pkg/config"
	servicecommon "github.com/BookManagementSystem/backend/pkg/service/common"
	"github.com/BookManagementSystem/backend/pkg/store"
	storecommon "github.com/BookManagementSystem/backend/pkg/store/common"
	"github.com/ant0ine/go-json-rest/rest"
)

type BooksService struct {
	lg *slog.Logger

	books books.Books
	store store.BookStore
	api   *http.Server

	token string
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
		api: &http.Server{
			Addr: config.Address,
		},
		token: config.Token,
	}, nil
}

func (s *BooksService) Listen() error {
	api := rest.NewApi()
	api.Use(CORSMiddleware())
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/books", s.GetAllBooks),
		rest.Get("/book:isbn", s.GetBook),
		rest.Get("/books/search:title", s.SearchBook),
		rest.Post("/put:isbn", s.Put),
		rest.Delete("/book:isbn", s.Delete),
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
	s.lg.Info("Start Listen Service", slog.String("address", s.api.Addr))
	if err := s.api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}

	s.lg.Info("Close Listen Service")
	return nil
}

func (s *BooksService) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.api.Shutdown(ctx); err != nil {
		s.lg.Error("failed to shutdown api", slog.String("err", err.Error()))
		return fmt.Errorf("failed to shutdown api: %w", err)
	}
	s.lg.Info("Shutdown BooksService")
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

func (s *BooksService) Put(w rest.ResponseWriter, r *rest.Request) {
	if err := s.Authorization(r); err != nil {
		rest.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isbn := strings.ReplaceAll(r.PathParam("isbn"), ":", "")

	info, err := s.books.GetInfo(isbn)
	if err != nil {
		s.lg.Error("failed to get info", slog.String("err", err.Error()))
		return
	}
	s.lg.Info("Get book info", slog.String("title", info.Title), slog.String("isbn", info.ISBN))
	if err := s.store.Put(*info); err != nil {
		s.lg.Error("failed to put info in store", slog.String("err", err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *BooksService) Authorization(r *rest.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader != s.token {
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func (s *BooksService) GetBook(w rest.ResponseWriter, r *rest.Request) {
	if err := s.Authorization(r); err != nil {
		rest.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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
	if err := s.Authorization(r); err != nil {
		rest.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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
	if err := s.Authorization(r); err != nil {
		rest.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

func (s *BooksService) Delete(w rest.ResponseWriter, r *rest.Request) {
	if err := s.Authorization(r); err != nil {
		rest.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	isbn := strings.ReplaceAll(r.PathParam("isbn"), ":", "")

	if err := s.store.Del(isbn); err != nil {
		s.lg.Error("internal server error", slog.String("err", err.Error()))
		rest.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
