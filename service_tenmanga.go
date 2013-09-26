package mangadownloader

import (
	//"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

const (
	serviceTenMangaDomain = "www.tenmanga.com"
)

var (
	serviceTenMangaRegexpIdentifyManga, _   = regexp.Compile("^/book/.+$")
	serviceTenMangaRegexpIdentifyChapter, _ = regexp.Compile("^/chapter/.+$")
)

type TenMangaService struct {
	Md *MangaDownloader
}

func (service *TenMangaService) Supports(u *url.URL) bool {
	return u.Host == serviceTenMangaDomain
}

func (service *TenMangaService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	if serviceTenMangaRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	if serviceTenMangaRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
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
