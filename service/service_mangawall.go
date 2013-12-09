package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

var (
	serviceMangaWallHosts = []string{
		"mangawall.com",
		"www.mangawall.com",
	}

	serviceMangaWallUrlBase *url.URL

	serviceMangaWallHtmlSelectorMangaName, _          = selector.Selector("meta[name=og:title]")
	serviceMangaWallHtmlSelectorMangaChapters, _      = selector.Selector(".chapterlistfull a")
	serviceMangaWallHtmlSelectorChapterPagesSelect, _ = selector.Selector(".pageselect")
	serviceMangaWallHtmlSelectorChapterPagesOption, _ = selector.Selector("option")
	serviceMangaWallHtmlSelectorPageImage, _          = selector.Selector(".scan")

	serviceMangaWallRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z\\-]+/?$")
	serviceMangaWallRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z\\-]+/.+$")
	serviceMangaWallRegexpChapterName, _     = regexp.Compile("^/manga/[0-9a-z\\-]+/([0-9\\.\\-]+).*$")
	serviceMangaWallRegexpPageBaseUrlPath, _ = regexp.Compile("^(/manga/[0-9a-z\\-]+/[0-9\\.\\-]+).*$")
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
	matches := serviceMangaWallRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name := matches[1]

	return name, nil
}

func (service *MangaWallService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	selectNodes := serviceMangaWallHtmlSelectorChapterPagesSelect.Find(rootNode)
	if len(selectNodes) != 2 {
		return nil, errors.New("Select node not found")
	}
	selectNode := selectNodes[0]
	optionNodes := serviceMangaWallHtmlSelectorChapterPagesOption.Find(selectNode)

	matches := serviceMangaWallRegexpPageBaseUrlPath.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return nil, errors.New("Invalid path format")
	}
	pageBaseUrlPath := matches[1]

	pageBaseUrl := urlCopy(serviceMangaWallUrlBase)
	pageBaseUrl.Path = pageBaseUrlPath

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
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaWallHtmlSelectorPageImage.Find(rootNode)
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
