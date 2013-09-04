package mangadownloader

import (
	"net/url"
)

type Service interface {
	Mangas() ([]*Manga, error)
	MangaName(*Manga) (string, error)
	MangaChapters(*Manga) ([]*Chapter, error)
	ChapterPages(*Chapter) ([]*Page, error)
	PageIndex(*Page) (uint, error)
	PageImage(*Page) (*Image, error)
	Supports(*url.URL) bool
	Identify(*url.URL) (interface{}, error)
}
