package db

import (
	"fmt"
	"log/slog"

	bookscommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/common"
	storeconfig "github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/config"
	"github.com/nyahahanoha/BookManagementSystem/backend/pkg/store/db/mysql"
)

type DBStore interface {
	Init() error
	Put(book bookscommon.Info) error
	Get(isbn string) (bookscommon.Info, error)
	GetAll() ([]bookscommon.Info, error)
	Search(title string) ([]bookscommon.Info, error)
	Delete(isbn string) error

	Rename(isbn, title string) error

	Close() error
}

func NewDBStore(lg *slog.Logger, config storeconfig.DBConfig) (DBStore, error) {
	switch config.Kind {
	case storeconfig.MySQL:
		db, err := mysql.NewMySQL(lg, config.MySQLConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect mysql: %w", err)
		}
		if err := db.Init(); err != nil {
			return nil, fmt.Errorf("failed to init database: %w", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("failed to connect db: invalid kind")
	}
}
