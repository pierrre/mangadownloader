package mangadownloader

import (
	"errors"
	"net/url"
	"regexp"
)

const (
	serviceMangaWallHtmlSelectorMangaName          = "meta[name='og:title']"
	serviceMangaWallHtmlSelectorMangaChapters      = ".chapterlistfull a"
	serviceMangaWallHtmlSelectorChapterPagesSelect = ".pageselect"
	serviceMangaWallHtmlSelectorChapterPagesOption = "option"
	serviceMangaWallHtmlSelectorPageImage          = ".scan"
)

var (
	serviceMangaWallHosts = []string{
		"mangawall.com",
		"www.mangawall.com",
	}

	serviceMangaWallUrlBase *url.URL

	serviceMangaWallRegexpIdentifyManga   = regexp.MustCompile("^/manga/[0-9a-z\\-]+/?$")
	serviceMangaWallRegexpIdentifyChapter = regexp.MustCompile("^/manga/[0-9a-z\\-]+/.+$")
	serviceMangaWallRegexpChapterName     = regexp.MustCompile("^/manga/[0-9a-z\\-]+/([0-9\\.\\-]+).*$")
	serviceMangaWallRegexpPageBaseUrlPath = regexp.MustCompile("^(/manga/[0-9a-z\\-]+/[0-9\\.\\-]+).*$")
)

func init() {
	serviceMangaWallUrlBase = new(url.URL)
	serviceMangaWallUrlBase.Scheme = "http"
	serviceMangaWallUrlBase.Host = serviceMangaWallHosts[0]
}

type MangaWallService struct {
	Md *MangaDownloader
}

func (service *MangaWallService) Supports(u *url.URL) bool {
	return stringSliceContains(serviceMangaWallHosts, u.Host)
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
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return "", err
	}

	metaOgTitleNodes := doc.Find(serviceMangaWallHtmlSelectorMangaName).Nodes
	if len(metaOgTitleNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	metaOgTitleNode := metaOgTitleNodes[0]
	name := htmlGetNodeAttribute(metaOgTitleNode, "content")

	return name, nil
}

func (service *MangaWallService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := doc.Find(serviceMangaWallHtmlSelectorMangaChapters).Nodes
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
	matches := serviceMangaWallRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name := matches[1]

	return name, nil
}

func (service *MangaWallService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	doc, err := service.Md.HttpGetHtmlDoc(chapter.Url)
	if err != nil {
		return nil, err
	}

	matches := serviceMangaWallRegexpPageBaseUrlPath.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return nil, errors.New("Invalid path format")
	}
	pageBaseUrlPath := matches[1]

	pageBaseUrl := urlCopy(serviceMangaWallUrlBase)
	pageBaseUrl.Path = pageBaseUrlPath

	optionNodes := doc.Find(serviceMangaWallHtmlSelectorChapterPagesSelect).First().Find(serviceMangaWallHtmlSelectorChapterPagesOption).Nodes
	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageUrl := urlCopy(pageBaseUrl)
		pageUrl.Path += "/" + htmlGetNodeAttribute(optionNode, "value")
		page := &Page{
			Url:     pageUrl,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaWallService) PageImageUrl(page *Page) (*url.URL, error) {
	doc, err := service.Md.HttpGetHtmlDoc(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := doc.Find(serviceMangaWallHtmlSelectorPageImage).Nodes
	if len(imgNodes) != 1 {
		return nil, errors.New("Image node not found")
	}
	imgNode := imgNodes[0]

	imageUrl, err := url.Parse(htmlGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageUrl, nil
}

func (service *MangaWallService) String() string {
	return "MangaWallService"
}
