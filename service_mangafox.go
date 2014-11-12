package mangadownloader

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

const (
	serviceMangaFoxHtmlSelectorMangaName      = "#series_info div.cover img"
	serviceMangaFoxHtmlSelectorMangaChapters1 = "#chapters ul.chlist li h3 a"
	serviceMangaFoxHtmlSelectorMangaChapters2 = "#chapters ul.chlist li h4 a"
	serviceMangaFoxHtmlSelectorChapterPages   = "#top_center_bar div.r option"
	serviceMangaFoxHtmlSelectorPageImage      = "#image"
)

var (
	serviceMangaFoxHosts = []string{
		"mangafox.me",
		"beta.mangafox.com",
	}

	serviceMangaFoxRegexpIdentifyManga   = regexp.MustCompile("^/manga/[0-9a-z_]+/?$")
	serviceMangaFoxRegexpIdentifyChapter = regexp.MustCompile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaFoxRegexpChapterName     = regexp.MustCompile("^.*/c(\\d+(\\.\\d+)?).*$")
	serviceMangaFoxRegexpPageBaseUrlPath = regexp.MustCompile("/?(\\d+\\.html)?$")
)

type MangaFoxService struct {
	Md *MangaDownloader
}

func (service *MangaFoxService) Supports(u *url.URL) bool {
	return stringSliceContains(serviceMangaFoxHosts, u.Host)
}

func (service *MangaFoxService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, fmt.Errorf("url '%s' not supported", u)
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

	return nil, fmt.Errorf("url '%s' unknown", u)
}

func (service *MangaFoxService) MangaName(manga *Manga) (string, error) {
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := doc.Find(serviceMangaFoxHtmlSelectorMangaName).Nodes
	if len(nameNodes) != 1 {
		return "", fmt.Errorf("html node '%s' (manga name) not found in '%s'", serviceMangaFoxHtmlSelectorMangaName, manga.Url)
	}
	nameNode := nameNodes[0]

	name := htmlGetNodeAttribute(nameNode, "alt")

	return name, nil
}

func (service *MangaFoxService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	doc, err := service.Md.HttpGetHtmlDoc(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := make([]*html.Node, 0)
	linkNodes = append(linkNodes, doc.Find(serviceMangaFoxHtmlSelectorMangaChapters1).Nodes...)
	linkNodes = append(linkNodes, doc.Find(serviceMangaFoxHtmlSelectorMangaChapters2).Nodes...)

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

func (service *MangaFoxService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaFoxRegexpChapterName.FindStringSubmatch(chapter.Url.Path)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (chapter name) not found in '%s'", serviceMangaFoxRegexpChapterName, chapter.Url)
	}
	name := matches[1]

	return name, nil
}

func (service *MangaFoxService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	doc, err := service.Md.HttpGetHtmlDoc(chapter.Url)
	if err != nil {
		return nil, err
	}

	basePageUrl := urlCopy(chapter.Url)
	basePageUrl.Path = serviceMangaFoxRegexpPageBaseUrlPath.ReplaceAllString(basePageUrl.Path, "")

	optionNodes := doc.Find(serviceMangaFoxHtmlSelectorChapterPages).Nodes
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
	doc, err := service.Md.HttpGetHtmlDoc(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := doc.Find(serviceMangaFoxHtmlSelectorPageImage).Nodes
	if len(imgNodes) != 1 {
		return nil, fmt.Errorf("html node '%s' (page image url) not found in '%s'", serviceMangaFoxHtmlSelectorPageImage, page.Url)
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
