package service

import (
	"github.com/matrixik/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
)

var (
	mangawall = &MangaWallService{}

	serviceMangaWallHTMLSelectorMangaName, _          = selector.Selector("meta[name=og:title]")
	serviceMangaWallHTMLSelectorMangaChapters, _      = selector.Selector(".chapterlistfull a")
	serviceMangaWallHTMLSelectorChapterPagesSelect, _ = selector.Selector(".pageselect")
	serviceMangaWallHTMLSelectorChapterPagesOption, _ = selector.Selector("option")
	serviceMangaWallHTMLSelectorPageImage, _          = selector.Selector(".scan")

	serviceMangaWallRegexpIdentifyManga, _   = regexp.Compile("^/manga/[0-9a-z\\-]+/?$")
	serviceMangaWallRegexpIdentifyChapter, _ = regexp.Compile("^/manga/[0-9a-z\\-]+/.+$")
	serviceMangaWallRegexpChapterName, _     = regexp.Compile("^/manga/[0-9a-z\\-]+/([0-9\\.\\-]+).*$")
	serviceMangaWallRegexpPageBaseURLPath, _ = regexp.Compile("^(/manga/[0-9a-z\\-]+/[0-9\\.\\-]+).*$")
)

type MangaWallService struct {
	ServiceCommon
}

func init() {
	mangawall.ServiceCommon = ServiceCommon{
		Hosts: []string{
			"mangawall.com",
			"www.mangawall.com",
		},
	}

	mangawall.URLBase = new(url.URL)
	mangawall.URLBase.Scheme = "http"
	mangawall.URLBase.Host = mangawall.Hosts[0]

	RegisterService("mangawall", mangawall)
}

func (service *MangaWallService) Supports(u *url.URL) bool {
	return utils.StringSliceContains(mangawall.Hosts, u.Host)
}

func (service *MangaWallService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	if serviceMangaWallRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			URL:     u,
			Service: service,
		}
		return chapter, nil
	}

	if serviceMangaWallRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			URL:     u,
			Service: service,
		}
		return manga, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *MangaWallService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	metaOgTitleNodes := serviceMangaWallHTMLSelectorMangaName.Find(rootNode)
	if len(metaOgTitleNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	metaOgTitleNode := metaOgTitleNodes[0]
	name := utils.HTMLGetNodeAttribute(metaOgTitleNode, "content")

	return name, nil
}

func (service *MangaWallService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaWallHTMLSelectorMangaChapters.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterURL := utils.URLCopy(mangawall.URLBase)
		chapterURL.Path = utils.HTMLGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			URL:     chapterURL,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *MangaWallService) ChapterName(chapter *Chapter) (string, error) {
	matches := serviceMangaWallRegexpChapterName.FindStringSubmatch(chapter.URL.Path)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name := matches[1]

	return name, nil
}

func (service *MangaWallService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	selectNodes := serviceMangaWallHTMLSelectorChapterPagesSelect.Find(rootNode)
	if len(selectNodes) != 2 {
		return nil, errors.New("Select node not found")
	}
	selectNode := selectNodes[0]
	optionNodes := serviceMangaWallHTMLSelectorChapterPagesOption.Find(selectNode)

	matches := serviceMangaWallRegexpPageBaseURLPath.FindStringSubmatch(chapter.URL.Path)
	if matches == nil {
		return nil, errors.New("Invalid path format")
	}
	pageBaseURLPath := matches[1]

	pageBaseURL := utils.URLCopy(mangawall.URLBase)
	pageBaseURL.Path = pageBaseURLPath

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageURL := utils.URLCopy(pageBaseURL)
		pageURL.Path += "/" + utils.HTMLGetNodeAttribute(optionNode, "value")
		page := &Page{
			URL:     pageURL,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaWallService) PageImageURL(page *Page) (*url.URL, error) {
	rootNode, err := utils.HTTPGetHTML(page.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaWallHTMLSelectorPageImage.Find(rootNode)
	if len(imgNodes) != 1 {
		return nil, errors.New("Image node not found")
	}
	imgNode := imgNodes[0]

	imageURL, err := url.Parse(utils.HTMLGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageURL, nil
}

func (service *MangaWallService) HTTPRetry() int {
	return service.httpRetry
}
func (service *MangaWallService) SetHTTPRetry(nr int) {
	service.httpRetry = nr
}

func (service *MangaWallService) String() string {
	return "MangaWallService"
}
