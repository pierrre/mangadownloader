package mangadownloader

import (
	//"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

const (
	serviceMangaWallDomain = "mangawall.com"
)

var (
	serviceMangaWallRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z\\-]+/?$")
	serviceMangaWallRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z\\-]+/.+$")
)

type MangaWallService struct {
	Md *MangaDownloader
}

func (service *MangaWallService) Supports(u *url.URL) bool {
	return u.Host == serviceMangaWallDomain
}

func (service *MangaWallService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	if serviceMangaWallRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaWallRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *MangaWallService) MangaName(manga *Manga) (string, error) {
	return "", errors.New("Not implemented")
}

func (service *MangaWallService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	return nil, errors.New("Not implemented")
}

func (service *MangaWallService) ChapterName(chapter *Chapter) (string, error) {
	return "", errors.New("Not implemented")
}

func (service *MangaWallService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	return nil, errors.New("Not implemented")
}

func (service *MangaWallService) PageImageUrl(page *Page) (*url.URL, error) {
	return nil, errors.New("Not implemented")
}

func (service *MangaWallService) String() string {
	return "MangaWallService"
}
