package googlebooks

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	bookscommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/common"
	booksconfig "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/config"
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
	if len(volumes.Items) == 0 {
		return nil, fmt.Errorf("not found book")
	}

	volume := volumes.Items[0]
	var otherTitle []string
	for _, item := range volumes.Items {
		volumeFullTitle := volume.VolumeInfo.Title + " " + volume.VolumeInfo.Subtitle
		itemFullTitle := item.VolumeInfo.Title + " " + item.VolumeInfo.Subtitle

		if strings.Contains(volumeFullTitle, itemFullTitle) || strings.Contains(itemFullTitle, volumeFullTitle) {
			if len(itemFullTitle) > len(volumeFullTitle) {
				volume = item
			}
			continue
		}

		otherTitle = append(otherTitle, itemFullTitle)

		itemDate, err := StringToDate(item.VolumeInfo.PublishedDate)
		if err != nil {
			continue
		}
		volumeDate, err := StringToDate(volume.VolumeInfo.PublishedDate)
		if err != nil {
			volume.VolumeInfo.PublishedDate = item.VolumeInfo.PublishedDate
			continue
		}

		if itemDate.After(volumeDate) {
			volume.VolumeInfo.PublishedDate = item.VolumeInfo.PublishedDate
			continue
		}
	}

	title := volume.VolumeInfo.Title + " " + volume.VolumeInfo.Subtitle
	if len(otherTitle) > 0 {
		title = fmt.Sprintf("%s / %s", title, strings.Join(otherTitle, " / "))
	}

	var desc string
	if volume.VolumeInfo.Description != "" {
		desc = volume.VolumeInfo.Description
	} else if volume.SearchInfo != nil {
		desc = volume.SearchInfo.TextSnippet
	} else {
		desc = bookscommon.NoDescription
	}

	book := &bookscommon.Info{
		ISBN:        isbn,
		Title:       title,
		Authors:     volume.VolumeInfo.Authors,
		Description: desc,
		Language:    StringToLanguage(volume.VolumeInfo.Language),
	}

	date, err := StringToDate(volume.VolumeInfo.PublishedDate)
	if err == nil {
		book.Publishdate = date
	}

	var u *url.URL
	if volume.VolumeInfo.ImageLinks != nil {
		u, err = url.Parse(volume.VolumeInfo.ImageLinks.Thumbnail)
		if err == nil {
			book.Image.Source = *u
		}
	}

	return book, nil
}

func StringToDate(s string) (time.Time, error) {
	var date time.Time
	if s == "" {
		return date, fmt.Errorf("date is empty")
	}
	if len(s) == 4 {
		s = s + "-01"
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
