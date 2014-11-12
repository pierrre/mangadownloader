package mangadownloader

import (
	"fmt"
	"net/url"
	"regexp"
)

const (
	serviceMangaHereHtmlSelectorMangaName     = ".detail_list .title h3"
	serviceMangaHereHtmlSelectorMangaChapters = ".detail_list a"
	serviceMangaHereHtmlSelectorChapterPages  = ".readpage_top .right option"
	serviceMangaHereHtmlSelectorPageImage     = "#image"
)

var (
	serviceMangaHereHosts = []string{
		"www.mangahere.com",
		"mangahere.com",
	}

	serviceMangaHereRegexpIdentifyManga   = regexp.MustCompile("^/manga/[0-9a-z_]+/?$")
	serviceMangaHereRegexpIdentifyChapter = regexp.MustCompile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaHereRegexpMangaName       = regexp.MustCompile("^Read (.*) Online$")
	serviceMangaHereRegexpChapterName     = regexp.MustCompile("^.*/c(\\d+(\\.\\d+)?).*$")
)

type MangaHereService struct {
	Md *MangaDownloader
}

func (service *MangaHereService) Supports(u *url.URL) bool {
	return stringSliceContains(serviceMangaHereHosts, u.Host)
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
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := doc.Find(serviceMangaHereHtmlSelectorMangaName).Nodes
	if len(nameNodes) != 1 {
		return "", fmt.Errorf("html node '%s' (manga name) not found in '%s'", serviceMangaHereHtmlSelectorMangaName, manga.Url)
	}
	nameNode := nameNodes[0]

	name, err := htmlGetNodeText(nameNode)
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
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := doc.Find(serviceMangaHereHtmlSelectorMangaChapters).Nodes
	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl, err := url.Parse(htmlGetNodeAttribute(linkNode, "href"))
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
	doc, err := service.Md.HttpGetHtmlDoc(chapter.Url)
	if err != nil {
		return nil, err
	}

	optionNodes := doc.Find(serviceMangaHereHtmlSelectorChapterPages).Nodes
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
	doc, err := service.Md.HttpGetHtmlDoc(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := doc.Find(serviceMangaHereHtmlSelectorPageImage).Nodes
	if len(imgNodes) != 1 {
		return nil, fmt.Errorf("html node '%s' (page image url) not found in '%s'", serviceMangaHereHtmlSelectorPageImage, page.Url)
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
