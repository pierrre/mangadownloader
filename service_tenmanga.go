package mangadownloader

import (
	//"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	//"regexp"
)

const (
	serviceTenMangaDomain = "www.tenmanga.com"
)

type TenMangaService struct {
	Md *MangaDownloader
}

func (service *TenMangaService) Supports(u *url.URL) bool {
	return u.Host == serviceTenMangaDomain
}

func (service *TenMangaService) Identify(u *url.URL) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (service *TenMangaService) MangaName(manga *Manga) (string, error) {
	return "", errors.New("Not implemented")
}

func (service *TenMangaService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	return nil, errors.New("Not implemented")
}

func (service *TenMangaService) ChapterName(chapter *Chapter) (string, error) {
	return "", errors.New("Not implemented")
}

func (service *TenMangaService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	return nil, errors.New("Not implemented")
}

func (service *TenMangaService) PageImageUrl(page *Page) (*url.URL, error) {
	return nil, errors.New("Not implemented")
}

func (service *TenMangaService) String() string {
	return "TenMangaService"
}
