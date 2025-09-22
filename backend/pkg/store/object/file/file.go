package filestore

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	storeconfig "github.com/BookManagementSystem/pkg/store/config"
)

type FileStore struct {
	lg     *slog.Logger
	prefix string
}

func NewFileStore(lg *slog.Logger, config storeconfig.FileConfig) (*FileStore, error) {
	return &FileStore{
		prefix: config.Prefix,
		lg:     lg.With(slog.String("Package", "filesystem")),
	}, nil
}

func (s *FileStore) Put(url url.URL, isbn string) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.lg.Error("failed to close responce")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	ct := resp.Header.Get("Content-Type")
	var ext string
	switch ct {
	case "image/jpeg":
		ext = ".jpeg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	default:
		ext = ""
	}
	path := path.Join(s.prefix, isbn) + ext
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *FileStore) Get(isbn string) (string, error) {
	entries, err := os.ReadDir(s.prefix)
	if err != nil {
		return "", fmt.Errorf("failed to read dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), isbn) {
			return entry.Name(), nil
		}
	}
	return "", nil
}
