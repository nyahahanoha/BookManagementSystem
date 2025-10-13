package ndlbooks

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	bookscommon "github.com/nyahahanoha/BookManagementSystem/backend/pkg/books/common"
)

// RSS構造体 (必要最小限)
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title    string   `xml:"title"`
	Authors  []string `xml:"creator"`
	PubDate  string   `xml:"pubDate"`
	Language string   `xml:"publicationPlace"`
	Volume   string   `xml:"volume"`
}

func StringToDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC1123Z, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC1123, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("failed to convert date: %s", s)
}

// ItemからBookInfoへ変換
func ToBookInfo(item Item) bookscommon.Info {
	fullTitle := item.Title + " " + item.Volume

	date, err := StringToDate(item.PubDate)
	if err != nil {
		date = time.Time{}
	}

	return bookscommon.Info{
		Title:       fullTitle,
		Authors:     item.Authors,
		Publishdate: date,
	}
}

type NDL struct{}

func NewNDL() (*NDL, error) {
	return &NDL{}, nil
}

func (s *NDL) Close() error {
	return nil
}

func (s *NDL) GetInfo(isbn string) (*bookscommon.Info, error) {
	u := fmt.Sprintf("https://ndlsearch.ndl.go.jp/api/opensearch?isbn=%s", isbn)
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	if len(rss.Channel.Items) == 0 {
		return nil, fmt.Errorf("no book found for ISBN: %s", isbn)
	}

	item := rss.Channel.Items[0]
	info := ToBookInfo(item)
	for _, item := range rss.Channel.Items {
		fullTitle := item.Title + " " + item.Volume

		titlePre := strings.Split(info.Title, " ")
		titleNext := strings.Split(fullTitle, " ")

		unique := make(map[string]struct{})
		for _, w := range titleNext {
			unique[w] = struct{}{}
		}
		for _, w := range titlePre {
			unique[w] = struct{}{}
		}
		var merged []string
		for w := range unique {
			merged = append(merged, w)
		}
		info.Title = strings.Join(merged, " ")
	}

	info.ISBN = isbn
	info.Language = bookscommon.JP
	info.Description = bookscommon.NoDescription

	imgUrl, err := url.Parse("https://ndlsearch.ndl.go.jp/thumbnail/" + isbn + ".jpg")
	if err == nil {
		imgResp, err := http.Head(imgUrl.String())
		if err != nil {
			return &info, nil
		}
		defer imgResp.Body.Close()
		if imgResp.StatusCode == http.StatusOK {
			info.Image.Source = *imgUrl
		}
	}

	return &info, nil
}
