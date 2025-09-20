package bookscommon

import (
	"time"
)

type Language int32

const (
	UNKOWN Language = iota
	JP
	EN
)

type Info struct {
	Title       string
	Authoers    []string
	Description string
	Publishdate time.Time
	Language    Language
}
