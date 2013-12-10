package service

import (
	"net/url"
)

type Manga struct {
	URL     *url.URL
	Service ServiceHandler
}

func (manga *Manga) Name() (string, error) {
	return manga.Service.MangaName(manga)
}

func (manga *Manga) Chapters() ([]*Chapter, error) {
	return manga.Service.MangaChapters(manga)
}
