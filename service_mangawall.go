package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

const (
	serviceMangaWallDomain = "mangawall.com"
)

var (
	serviceMangaWallUrlBase *url.URL

	serviceMangaWallHtmlSelectorMangaName, _     = selector.Selector("meta[name=og:title]")
	serviceMangaWallHtmlSelectorMangaChapters, _ = selector.Selector(".chapterlistfull a")

	serviceMangaWallRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z\\-]+/?$")
	serviceMangaWallRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z\\-]+/.+$")
)

func init() {
	serviceMangaWallUrlBase = new(url.URL)
	serviceMangaWallUrlBase.Scheme = "http"
	serviceMangaWallUrlBase.Host = serviceMangaWallDomain
}

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
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return "", err
	}

	metaOgTitleNodes := serviceMangaWallHtmlSelectorMangaName.Find(rootNode)
	if len(metaOgTitleNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	metaOgTitleNode := metaOgTitleNodes[0]
	name := htmlGetNodeAttribute(metaOgTitleNode, "content")

	return name, nil
}

func (service *MangaWallService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaWallHtmlSelectorMangaChapters.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl := urlCopy(serviceMangaWallUrlBase)
		chapterUrl.Path = htmlGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
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
