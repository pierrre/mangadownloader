package service

import (
	"github.com/matrixik/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

var (
	mangafox = &MangaFoxService{}

	serviceMangaFoxHTMLSelectorMangaName, _      = selector.Selector("#series_info div.cover img")
	serviceMangaFoxHTMLSelectorMangaChapters1, _ = selector.Selector("#chapters ul.chlist li h3 a")
	serviceMangaFoxHTMLSelectorMangaChapters2, _ = selector.Selector("#chapters ul.chlist li h4 a")
	serviceMangaFoxHTMLSelectorChapterPages, _   = selector.Selector("#top_center_bar div.r option")
	serviceMangaFoxHTMLSelectorPageImage, _      = selector.Selector("#image")

	serviceMangaFoxRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z_]+/?$")
	serviceMangaFoxRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z_]+/.+$")
	serviceMangaFoxRegexpChapterName, _     = regexp.Compile("^.*/c(\\d+(\\.\\d+)?).*$")
	serviceMangaFoxRegexpPageBaseURLPath, _ = regexp.Compile("/?(\\d+\\.html)?$")
)

type MangaFoxService struct {
	ServiceCommon
}

func init() {
	mangafox.ServiceCommon = ServiceCommon{
		Hosts: []string{
			"mangafox.me",
			"beta.mangafox.com",
		},
	}

	RegisterService("mangafox", mangafox)
}

func (service *MangaFoxService) Supports(u *url.URL) bool {
	return utils.StringSliceContains(mangafox.Hosts, u.Host)
}

func (service *MangaFoxService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, fmt.Errorf("url '%s' not supported", u)
	}

	if serviceMangaFoxRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			URL:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaFoxRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			URL:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, fmt.Errorf("url '%s' unknown", u)
}

func (service *MangaFoxService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaFoxHTMLSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", fmt.Errorf("html node '%s' (manga name) not found in '%s'",
			serviceMangaFoxHTMLSelectorMangaName, manga.URL)
	}
	nameNode := nameNodes[0]

	name := utils.HTMLGetNodeAttribute(nameNode, "alt")

	return name, nil
}

func (service *MangaFoxService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := make([]*html.Node, 0)
	linkNodes = append(linkNodes, serviceMangaFoxHTMLSelectorMangaChapters1.Find(rootNode)...)
	linkNodes = append(linkNodes, serviceMangaFoxHTMLSelectorMangaChapters2.Find(rootNode)...)

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

func (service *MangaFoxService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaFoxRegexpChapterName.FindStringSubmatch(chapter.URL.Path)
	if matches == nil {
		return "", fmt.Errorf("regexp '%s' (chapter name) not found in '%s'",
			serviceMangaFoxRegexpChapterName, chapter.URL)
	}
	name := matches[1]

	return name, nil
}

func (service *MangaFoxService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	basePageURL := utils.URLCopy(chapter.URL)
	basePageURL.Path = serviceMangaFoxRegexpPageBaseURLPath.ReplaceAllString(basePageURL.Path, "")

	optionNodes := serviceMangaFoxHTMLSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageNumberString := utils.HTMLGetNodeAttribute(optionNode, "value")
		pageNumber, err := strconv.Atoi(pageNumberString)
		if err != nil {
			return nil, err
		}

		if pageNumber <= 0 {
			continue
		}

		pageURL := utils.URLCopy(basePageURL)
		pageURL.Path += fmt.Sprintf("/%d.html", pageNumber)

		page := &Page{
			URL:     pageURL,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaFoxService) PageImageURL(page *Page) (*url.URL, error) {
	rootNode, err := utils.HTTPGetHTML(page.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaFoxHTMLSelectorPageImage.Find(rootNode)
	if len(imgNodes) != 1 {
		return nil, fmt.Errorf("html node '%s' (page image url) not found in '%s'",
			serviceMangaFoxHTMLSelectorPageImage, page.URL)
	}
	imgNode := imgNodes[0]

	imageURL, err := url.Parse(utils.HTMLGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageURL, nil
}

func (service *MangaFoxService) HTTPRetry() int {
	return service.httpRetry
}
func (service *MangaFoxService) SetHTTPRetry(nr int) {
	service.httpRetry = nr
}

func (service *MangaFoxService) String() string {
	return "MangaFoxService"
}
