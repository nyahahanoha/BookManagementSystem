package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	bookscommon "github.com/BookManagementSystem/backend/pkg/books/common"
	storecommon "github.com/BookManagementSystem/backend/pkg/store/common"
	storeconfig "github.com/BookManagementSystem/backend/pkg/store/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	lg *slog.Logger
	db *sql.DB
}

func NewMySQL(lg *slog.Logger, config storeconfig.MySQLConfig) (*MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true",
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
		title varchar(200), 
		description varchar(2000),
		publishdate date,
		language varchar(8),
		image varchar(200),
		deleted boolean DEFAULT false
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS authors(
		id int AUTO_INCREMENT PRIMARY KEY,
		isbn varchar(14),
		author varchar(200),
		deleted boolean DEFAULT false,
		UNIQUE KEY isbn_author (isbn, author)
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

	var pubDate interface{}
	if book.Publishdate.IsZero() {
		pubDate = nil
	} else {
		pubDate = book.Publishdate
	}

	_, err = tx.Exec(`INSERT INTO books(
		isbn,
		title,
		description,
		publishdate,
		language,
		image
	) VALUES (?, ?, ?, ?, ?, ?)
	 ON DUPLICATE KEY UPDATE
	  title = VALUES(title),
    description = VALUES(description),
    publishdate = VALUES(publishdate),
    language = VALUES(language),
    image = VALUES(image),
		deleted = false
	 `,
		book.ISBN,
		book.Title,
		book.Description,
		pubDate,
		book.Language.String(),
		book.Image.Source.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	for _, author := range book.Authors {
		_, err = tx.Exec(`INSERT INTO authors(
			isbn,
			author
		) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
		  isbn = VALUES(isbn),
 	    author = VALUES(author),
			deleted = false
		 `,
			book.ISBN,
			author,
		)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return nil
}

func (s *MySQL) Get(isbn string) (bookscommon.Info, error) {
	row, err := s.db.Query(`SELECT 
        isbn,
        title,
        description,
        publishdate,
        language,
        image
        FROM books WHERE isbn = ? AND deleted = false`, isbn)
	if err != nil {
		return bookscommon.Info{}, fmt.Errorf("failed to execute query: %w", err)
	}
	books, err := s.rowConvertInfo(row)
	if err != nil {
		return bookscommon.Info{}, fmt.Errorf("failed to convert info: %w", err)
	}
	if len(books) > 0 {
		return books[0], nil
	}
	return bookscommon.Info{}, storecommon.ErrNotFoundBook
}

func (s *MySQL) GetAll() ([]bookscommon.Info, error) {
	rows, err := s.db.Query(`SELECT 
        isbn,
        title,
        description,
        publishdate,
        language,
        image FROM books WHERE deleted = false`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	books, err := s.rowConvertInfo(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to convert info: %w", err)
	}

	return books, nil
}

func (s *MySQL) Search(title string) ([]bookscommon.Info, error) {
	rows, err := s.db.Query(`SELECT 
        isbn,
        title,
        description,
        publishdate,
        language,
        image
				FROM books WHERE title LIKE ? AND deleted = false`, "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	books, err := s.rowConvertInfo(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to convert info: %w", err)
	}

	return books, nil
}

func (s *MySQL) rowConvertInfo(rows *sql.Rows) ([]bookscommon.Info, error) {
	var books []bookscommon.Info
	for rows.Next() {
		var book bookscommon.Info
		var langStr, imgStr string
		var pubDate sql.NullTime
		err := rows.Scan(
			&book.ISBN,
			&book.Title,
			&book.Description,
			&pubDate,
			&langStr,
			&imgStr,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan book row: %w", err)
		}

		if pubDate.Valid {
			book.Publishdate = pubDate.Time
		} else {
			book.Publishdate = time.Time{}
		}

		book.Language, err = bookscommon.LanguageString(langStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get language: %w", err)
		}
		imgurl, err := url.Parse(imgStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get image url: %w", err)
		}
		book.Image.Source = *imgurl

		rows, err := s.db.Query(`SELECT author FROM authors WHERE isbn = ?`, book.ISBN)
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

		books = append(books, book)
	}
	return books, nil
}

func (s *MySQL) Delete(isbn string) error {
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

	if err := s.db.QueryRow(`UPDATE books SET deleted = true WHERE isbn = ?`, isbn); err.Err() != nil {
		return fmt.Errorf("failed to execute query: %w", err.Err())
	}
	if err := s.db.QueryRow(`UPDATE authors SET deleted = true WHERE isbn = ?`, isbn); err.Err() != nil {
		return fmt.Errorf("failed to execute query: %w", err.Err())
	}
	return nil
}
