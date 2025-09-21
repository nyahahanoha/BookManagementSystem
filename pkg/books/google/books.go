package googlebooks

import (
	"context"
	"fmt"
	"net/url"
	"time"

	bookscommon "github.com/BookManagementSystem/pkg/books/common"
	booksconfig "github.com/BookManagementSystem/pkg/books/config"
	api "google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

type GoogleBooks struct {
	svc *api.Service

	closer context.CancelFunc
}

func NewGoogleBooks(config booksconfig.GoogleBooksConfig) (*GoogleBooks, error) {
	ctx, cancel := context.WithCancel(context.Background())
	svc, err := api.NewService(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	return &GoogleBooks{
		svc:    svc,
		closer: cancel,
	}, nil
}

func (s *GoogleBooks) Close() error {
	s.closer()
	return nil
}

func (s *GoogleBooks) GetInfo(isbn string) (*bookscommon.Info, error) {
	volumes, err := s.svc.Volumes.List("isbn:" + isbn).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	if volumes.TotalItems > 0 {
		volume := volumes.Items[0]

		book := &bookscommon.Info{
			ISBN:        isbn,
			Title:       volume.VolumeInfo.Title,
			Authoers:    volume.VolumeInfo.Authors,
			Description: volume.VolumeInfo.Description,
			Language:    StringToLanguage(volume.VolumeInfo.Language),
		}

		date, err := StringToDate(volume.VolumeInfo.PublishedDate)
		if err == nil {
			book.Publishdate = date
		}

		url, err := url.Parse(volume.VolumeInfo.ImageLinks.Thumbnail)
		if err == nil {
			book.Image = *url
		}

		return book, nil

	} else {
		return nil, fmt.Errorf("failed to request: invalid isbn")
	}
}

func StringToDate(s string) (time.Time, error) {
	var date time.Time
	if s == "" {
		return date, fmt.Errorf("date is empty")
	}
	shortForm := "2006-01"
	date, err := time.Parse(shortForm, s)
	if err != nil {
		return date, fmt.Errorf("failed to convert date: %w", err)
	}
	return date, nil
}

func StringToLanguage(s string) bookscommon.Language {
	switch s {
	case "en":
		return bookscommon.EN
	case "ja":
		return bookscommon.JP
	default:
		return bookscommon.UNKOWN
	}
}
