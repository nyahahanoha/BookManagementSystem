package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	storeconfig "github.com/BookManagementSystem/pkg/store/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	lg *slog.Logger
	db *sql.DB
}

func NewMySQL(lg *slog.Logger, config storeconfig.MySQLConfig) (*MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		config.User,
		config.Password,
		config.IPAddress,
		config.Port,
		config.Database),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect mysql: %w", err)
	}
	return &MySQL{
		db: db,
		lg: lg.With(slog.String("Package", "mysql")),
	}, nil
}

func (s *MySQL) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (s *MySQL) Init() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS books(
		isbn varchar(14) PRIMARY KEY, 
		title varcher(200), 
		description varcher(200),
		publishdate date,
		language varcher(8)
		image varcher(200)
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS authors(
		id int AUTO_INCREMENT PRIMARY KEY,
		isbn varcher(14),
		author varcher(200),
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (s *MySQL) Put(book bookscommon.Info) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				s.lg.Error("failed to rollback at transaction", slog.String("err", err.Error()))
			}
			s.lg.Error("failed to reconver", slog.Any("p", p))
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				s.lg.Error("failed to rollback at transaction", slog.String("err", err.Error()))
			}
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.Exec(`INSERT INTO books(
		isbn,
		title,
		description,
		publishdate,
		language,
		image,
	) VALUES (?, ?, ?, ?, ?, ?)`,
		book.ISBN,
		book.Title,
		book.Description,
		book.Publishdate,
		book.Language.String(),
		book.Image.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	for _, author := range book.Authors {
		_, err = tx.Exec(`INSERT INTO authors(
			isbn,
			author,
		) VALUES (?, ?, ?, ?, ?, ?)`,
			book.ISBN,
			author,
		)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return nil
}

func (s *MySQL) Get(isbn string) (*bookscommon.Info, error) {
	row := s.db.QueryRow(`SELECT 
        isbn,
        title,
        description,
        publishdate,
        language,
        image
        FROM books WHERE isbn = ?`, isbn)

	var book bookscommon.Info
	var langStr, imgStr string
	err := row.Scan(
		&book.ISBN,
		&book.Title,
		&book.Description,
		&book.Publishdate,
		&langStr,
		&imgStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan book row: %w", err)
	}

	book.Language, err = bookscommon.LanguageString(langStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get language: %w", err)
	}
	imgurl, err := url.Parse(imgStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get image url: %w", err)
	}
	book.Image = *imgurl

	rows, err := s.db.Query(`SELECT author FROM authors WHERE isbn = ?`, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to query authors: %w", err)
	}
	defer func() {
		if err != rows.Close() {
			s.lg.Error("failed to close query result", slog.String("err", err.Error()))
		}
	}()

	var authors []string
	for rows.Next() {
		var author string
		if err := rows.Scan(&author); err != nil {
			return nil, fmt.Errorf("failed to scan author row: %w", err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("authors rows iteration error: %w", err)
	}

	book.Authors = authors

	return &book, nil
}
