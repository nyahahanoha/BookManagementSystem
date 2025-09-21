package bookscommon

import (
	"net/url"
	"time"
)

//go:generate go run github.com/dmarkham/enumer -type=Language
type Language int32

const (
	UNKOWN Language = iota
	JP
	EN
)

type Info struct {
	ISBN        string
	Title       string
	Authors     []string
	Description string
	Publishdate time.Time
	Language    Language
	Image       Image
}

type Image struct {
	Source url.URL
	Path   string
}
