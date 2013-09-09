package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
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

	serviceMangaFoxHtmlSelectorIdentifyManga, _   = selector.Selector("#chapters")
	serviceMangaFoxHtmlSelectorIdentifyChapter, _ = selector.Selector("#top_chapter_list")
)

func init() {
	serviceMangaFoxUrlBase = new(url.URL)
	serviceMangaFoxUrlBase.Scheme = "http"
	serviceMangaFoxUrlBase.Host = serviceMangaFoxDomain

	serviceMangaFoxUrlMangas = urlCopy(serviceMangaFoxUrlBase)
	serviceMangaFoxUrlMangas.Path = serviceMangaFoxPathMangas
}

type MangaFoxService struct {
	Md *MangaDownloader
}

func (service *MangaFoxService) Supports(u *url.URL) bool {
	return u.Host == serviceMangaFoxDomain
}

func (service *MangaFoxService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	rootNode, err := service.Md.HttpGetHtml(u)
	if err != nil {
		return nil, err
	}

	identifyMangaNodes := serviceMangaFoxHtmlSelectorIdentifyManga.Find(rootNode)
	if len(identifyMangaNodes) == 1 {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	identifyChapterNodes := serviceMangaFoxHtmlSelectorIdentifyChapter.Find(rootNode)
	if len(identifyChapterNodes) == 1 {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
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
