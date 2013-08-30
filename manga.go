package mangadownloader

import (
	"net/url"
)

type Manga struct {
	Url     *url.URL
	Service Service
}

func (manga *Manga) Chapters() ([]*Chapter, error) {
	return manga.Service.Chapters(manga)
}
