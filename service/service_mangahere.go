package service

import (
	"github.com/pierrre/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"fmt"
	"net/url"
	"regexp"
)

var (
	mangahere = &MangaHereService{}

	serviceMangaHereHTMLSelectorMangaName, _     = selector.Selector(".detail_list .title h3")
	serviceMangaHereHTMLSelectorMangaChapters, _ = selector.Selector(".detail_list a")
	serviceMangaHereHTMLSelectorChapterPages, _  = selector.Selector(".readpage_top .right option")
	serviceMangaHereHTMLSelectorPageImage, _     = selector.Selector("#image")

	serviceMangaHereRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z_]+/?$")
	serviceMangaHereRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaHereRegexpMangaName, _       = regexp.Compile("^Read (.*) Online$")
	serviceMangaHereRegexpChapterName, _     = regexp.Compile("^.*/c(\\d+(\\.\\d+)?).*$")
)

type MangaHereService struct {
	ServiceCommon
}

func init() {
	mangahere.ServiceCommon = ServiceCommon{
		Hosts: []string{
			"www.mangahere.com",
			"mangahere.com",
		},
	}

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
			URL:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaHereRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			URL:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, fmt.Errorf("url '%s' unknown", u)
}

func (service *MangaHereService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaHereHTMLSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", fmt.Errorf("html node '%s' (manga name) not found in '%s'",
			serviceMangaHereHTMLSelectorMangaName, manga.URL)
	}
	nameNode := nameNodes[0]

	name, err := utils.HTMLGetNodeText(nameNode)
	if err != nil {
		return "", err
	}

	matches := serviceMangaHereRegexpMangaName.FindStringSubmatch(name)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (manga name) not found in '%s'",
			serviceMangaHereRegexpMangaName, name)
	}
	name = matches[1]

	return name, nil
}

func (service *MangaHereService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaHereHTMLSelectorMangaChapters.Find(rootNode)
	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterURL, err := url.Parse(utils.HTMLGetNodeAttribute(linkNode, "href"))
		if err != nil {
			return nil, err
		}
		chapter := &Chapter{
			URL:     chapterURL,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}
	chapters = chapterSliceReverse(chapters)

	return chapters, nil
}

func (service *MangaHereService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaHereRegexpChapterName.FindStringSubmatch(chapter.URL.Path)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (chapter name) not found in '%s'",
			serviceMangaHereRegexpChapterName, chapter.URL)
	}
	name := matches[1]

	return name, nil
}

func (service *MangaHereService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaHereHTMLSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageURL, err := url.Parse(utils.HTMLGetNodeAttribute(optionNode, "value"))
		if err != nil {
			return nil, err
		}
		page := &Page{
			URL:     pageURL,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaHereService) PageImageURL(page *Page) (*url.URL, error) {
	rootNode, err := utils.HTTPGetHTML(page.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaHereHTMLSelectorPageImage.Find(rootNode)
	if len(imgNodes) != 1 {
		return nil, fmt.Errorf("html node '%s' (page image url) not found in '%s'",
			serviceMangaHereHTMLSelectorPageImage, page.URL)
	}
	imgNode := imgNodes[0]

	imageURL, err := url.Parse(utils.HTMLGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageURL, nil
}

func (service *MangaHereService) HTTPRetry() int {
	return service.httpRetry
}
func (service *MangaHereService) SetHTTPRetry(nr int) {
	service.httpRetry = nr
}

func (service *MangaHereService) String() string {
	return "MangaHereService"
}
