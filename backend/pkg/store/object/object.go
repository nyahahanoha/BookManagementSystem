package object

import (
	"fmt"
	"log/slog"
	"net/url"

	storeconfig "github.com/BookManagementSystem/pkg/store/config"
	filestore "github.com/BookManagementSystem/pkg/store/object/file"
)

type ObjectStore interface {
	Put(url url.URL, isbn string) error
	Get(isbn string) (string, error)
	Delete(isbn string) error
}

func NewObjectStore(lg *slog.Logger, config storeconfig.ObjectConfig) (ObjectStore, error) {
	switch config.Kind {
	case storeconfig.FileSystem:
		filesystem, err := filestore.NewFileStore(lg, config.FileConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create file store: %w", err)
		}
		return filesystem, nil
	default:
		return nil, fmt.Errorf("failed to connect object: invalid component")
	}
}
