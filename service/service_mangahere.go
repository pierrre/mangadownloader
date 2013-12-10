package service

import (
	"github.com/matrixik/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"fmt"
	"net/url"
	"regexp"
)

var (
	mangahere = &MangaHereService{
		Hosts: []string{
			"www.mangahere.com",
			"mangahere.com",
		},
	}

	serviceMangaHereHtmlSelectorMangaName, _     = selector.Selector(".detail_list .title h3")
	serviceMangaHereHtmlSelectorMangaChapters, _ = selector.Selector(".detail_list a")
	serviceMangaHereHtmlSelectorChapterPages, _  = selector.Selector(".readpage_top .right option")
	serviceMangaHereHtmlSelectorPageImage, _     = selector.Selector("#image")

	serviceMangaHereRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z_]+/?$")
	serviceMangaHereRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaHereRegexpMangaName, _       = regexp.Compile("^Read (.*) Online$")
	serviceMangaHereRegexpChapterName, _     = regexp.Compile("^.*/c(\\d+(\\.\\d+)?).*$")
)

type MangaHereService Service

func init() {
	RegisterService("mangahere", mangahere)
}

func (service *MangaHereService) Supports(u *url.URL) bool {
	return utils.StringSliceContains(mangahere.Hosts, u.Host)
}

func (service *MangaHereService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, fmt.Errorf("url '%s' not supported", u)
	}

	if serviceMangaHereRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaHereRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, fmt.Errorf("url '%s' unknown", u)
}

func (service *MangaHereService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HttpGetHtml(manga.Url, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaHereHtmlSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", fmt.Errorf("html node '%s' (manga name) not found in '%s'", serviceMangaHereHtmlSelectorMangaName, manga.Url)
	}
	nameNode := nameNodes[0]

	name, err := utils.HtmlGetNodeText(nameNode)
	if err != nil {
		return "", err
	}

	matches := serviceMangaHereRegexpMangaName.FindStringSubmatch(name)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (manga name) not found in '%s'", serviceMangaHereRegexpMangaName, name)
	}
	name = matches[1]

	return name, nil
}

func (service *MangaHereService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HttpGetHtml(manga.Url, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaHereHtmlSelectorMangaChapters.Find(rootNode)
	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl, err := url.Parse(utils.HtmlGetNodeAttribute(linkNode, "href"))
		if err != nil {
			return nil, err
		}
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}
	chapters = chapterSliceReverse(chapters)

	return chapters, nil
}

func (service *MangaHereService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaHereRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (chapter name) not found in '%s'", serviceMangaHereRegexpChapterName, chapter.Url)
	}
	name := matches[1]

	return name, nil
}

func (service *MangaHereService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HttpGetHtml(chapter.Url, service.httpRetry)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaHereHtmlSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageUrl, err := url.Parse(utils.HtmlGetNodeAttribute(optionNode, "value"))
		if err != nil {
			return nil, err
		}
		page := &Page{
			Url:     pageUrl,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaHereService) PageImageUrl(page *Page) (*url.URL, error) {
	rootNode, err := utils.HttpGetHtml(page.Url, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaHereHtmlSelectorPageImage.Find(rootNode)
	if len(imgNodes) != 1 {
		return nil, fmt.Errorf("html node '%s' (page image url) not found in '%s'", serviceMangaHereHtmlSelectorPageImage, page.Url)
	}
	imgNode := imgNodes[0]

	imageUrl, err := url.Parse(utils.HtmlGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageUrl, nil
}

func (service *MangaHereService) HttpRetry() int {
	return service.httpRetry
}
func (service *MangaHereService) SetHttpRetry(nr int) {
	service.httpRetry = nr
}

func (service *MangaHereService) String() string {
	return "MangaHereService"
}
