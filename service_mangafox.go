package mangadownloader

import (
	"errors"
	"net/url"
)

const (
	serviceMangaFoxDomain     = "mangafox.me"
	serviceMangaFoxPathMangas = "/manga"
)

var (
	serviceMangaFoxUrlBase   *url.URL
	serviceMangaFoxUrlMangas *url.URL
)

func init() {
	serviceMangaFoxUrlBase = new(url.URL)
	serviceMangaFoxUrlBase.Scheme = "http"
	serviceMangaFoxUrlBase.Host = serviceMangaFoxDomain

	serviceMangaFoxUrlMangas = urlCopy(serviceMangaFoxUrlBase)
	serviceMangaFoxUrlMangas.Path = serviceMangaFoxPathMangas
}

type MangaFoxService struct {
}

func (service *MangaFoxService) Supports(u *url.URL) bool {
	return u.Host == serviceMangaFoxDomain
}

func (service *MangaFoxService) Identify(u *url.URL) (interface{}, error) {
	//TODO
	return nil, errors.New("Not implemented")
}

func (service *MangaFoxService) Mangas() ([]*Manga, error) {
	//TODO
	return nil, errors.New("Not implemented")
}

func (service *MangaFoxService) MangaName(manga *Manga) (string, error) {
	//TODO
	return "", errors.New("Not implemented")
}

func (service *MangaFoxService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	//TODO
	return nil, errors.New("Not implemented")
}

func (service *MangaFoxService) ChapterName(chapter *Chapter) (string, error) {
	//TODO
	return "", errors.New("Not implemented")
}

func (service *MangaFoxService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	//TODO
	return nil, errors.New("Not implemented")
}

func (service *MangaFoxService) PageImageUrl(page *Page) (*url.URL, error) {
	//TODO
	return nil, errors.New("Not implemented")
}

func (service *MangaFoxService) String() string {
	return "MangaFoxService"
}
