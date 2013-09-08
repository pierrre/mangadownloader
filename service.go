package mangadownloader

import (
	"net/url"
)

type Service interface {
	Mangas() ([]*Manga, error)
	MangaName(*Manga) (string, error)
	MangaChapters(*Manga) ([]*Chapter, error)
	ChapterName(*Chapter) (string, error)
	ChapterPages(*Chapter) ([]*Page, error)
	PageImageUrl(*Page) (*url.URL, error)
	Supports(*url.URL) bool
	Identify(*url.URL) (interface{}, error)
}
