package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

var (
	serviceMangaHereHosts = []string{
		"www.mangahere.com",
		"mangahere.com",
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

type MangaHereService struct {
	Md *MangaDownloader
}

func (service *MangaHereService) Supports(u *url.URL) bool {
	return sliceStringContains(serviceMangaHereHosts, u.Host)
}

func (service *MangaHereService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
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

	return nil, errors.New("Unknown url")
}

func (service *MangaHereService) MangaName(manga *Manga) (string, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaHereHtmlSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data
	matches := serviceMangaHereRegexpMangaName.FindStringSubmatch(name)
	if matches == nil || len(matches) != 2 {
		return "", errors.New("Invalid name format")
	}
	name = matches[1]

	return name, nil
}

func (service *MangaHereService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaHereHtmlSelectorMangaChapters.Find(rootNode)
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

func (service *MangaHereService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaHereRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name := matches[1]

	return name, nil
}

func (service *MangaHereService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaHereHtmlSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageUrl, err := url.Parse(htmlGetNodeAttribute(optionNode, "value"))
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
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaHereHtmlSelectorPageImage.Find(rootNode)
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

func (service *MangaHereService) String() string {
	return "MangaHereService"
}
