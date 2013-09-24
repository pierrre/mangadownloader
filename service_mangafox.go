package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

const (
	serviceMangaFoxDomain = "mangafox.me"
)

var (
	serviceMangaFoxHtmlSelectorMangaName, _      = selector.Selector("#series_info div.cover img")
	serviceMangaFoxHtmlSelectorMangaChapters1, _ = selector.Selector("#chapters ul.chlist li h3 a")
	serviceMangaFoxHtmlSelectorMangaChapters2, _ = selector.Selector("#chapters ul.chlist li h4 a")
	serviceMangaFoxHtmlSelectorChapterPages, _   = selector.Selector("#top_center_bar div.r option")
	serviceMangaFoxHtmlSelectorPageImage, _      = selector.Selector("#image")

	serviceMangaFoxRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z_]+/?$")
	serviceMangaFoxRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaFoxRegexpChapterName, _     = regexp.Compile("^.*/c(\\d+(\\.\\d+)?).*$")
	serviceMangaFoxRegexpPageBaseUrlPath, _ = regexp.Compile("/?(\\d+\\.html)?$")
)

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

	if serviceMangaFoxRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaFoxRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *MangaFoxService) MangaName(manga *Manga) (string, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaFoxHtmlSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	name := htmlGetNodeAttribute(nameNode, "alt")

	return name, nil
}

func (service *MangaFoxService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := make([]*html.Node, 0)
	linkNodes = append(linkNodes, serviceMangaFoxHtmlSelectorMangaChapters1.Find(rootNode)...)
	linkNodes = append(linkNodes, serviceMangaFoxHtmlSelectorMangaChapters2.Find(rootNode)...)

	chaptersReversed := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl, err := url.Parse(htmlGetNodeAttribute(linkNode, "href"))
		if err != nil {
			return nil, err
		}
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chaptersReversed = append(chaptersReversed, chapter)
	}

	chapterCount := len(chaptersReversed)
	chapters := make([]*Chapter, 0, chapterCount)
	for i := chapterCount - 1; i >= 0; i-- {
		chapters = append(chapters, chaptersReversed[i])
	}

	return chapters, nil
}

func (service *MangaFoxService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaFoxRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name := matches[1]

	return name, nil
}

func (service *MangaFoxService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	basePageUrl := urlCopy(chapter.Url)
	basePageUrl.Path = serviceMangaFoxRegexpPageBaseUrlPath.ReplaceAllString(basePageUrl.Path, "")

	optionNodes := serviceMangaFoxHtmlSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageNumberString := htmlGetNodeAttribute(optionNode, "value")
		pageNumber, err := strconv.Atoi(pageNumberString)
		if err != nil {
			return nil, err
		}
		if pageNumber <= 0 {
			continue
		}
		pageUrl := urlCopy(basePageUrl)
		pageUrl.Path += fmt.Sprintf("/%d.html", pageNumber)
		page := &Page{
			Url:     pageUrl,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaFoxService) PageImageUrl(page *Page) (*url.URL, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaFoxHtmlSelectorPageImage.Find(rootNode)
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

func (service *MangaFoxService) String() string {
	return "MangaFoxService"
}
