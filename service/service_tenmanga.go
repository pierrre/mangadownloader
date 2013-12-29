package service

import (
	"github.com/matrixik/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var (
	tenmanga = &TenMangaService{
		Hosts: []string{
			"www.tenmanga.com",
			"tenmanga.com",
		},
	}

	serviceTenMangaHTMLSelectorMangaName, _        = selector.Selector(".postion .red")
	serviceTenMangaHTMLSelectorMangaChapters, _    = selector.Selector(".chapter_list td[align=left] a")
	serviceTenMangaHTMLSelectorChapterNameTitle, _ = selector.Selector("title")
	serviceTenMangaHTMLSelectorChapterNameManga, _ = selector.Selector(".postion a")
	serviceTenMangaHTMLSelectorChapterPages, _     = selector.Selector("#page option")
	serviceTenMangaHTMLSelectorPageImage, _        = selector.Selector("#comicpic")

	serviceTenMangaRegexpIdentifyManga, _   = regexp.Compile("^/book/.+$")
	serviceTenMangaRegexpIdentifyChapter, _ = regexp.Compile("^/chapter/.+$")
	serviceTenMangaRegexpChapterName, _     = regexp.Compile("^(.+) Page \\d+$")
)

func init() {
	tenmanga.URLBase = new(url.URL)
	tenmanga.URLBase.Scheme = "http"
	tenmanga.URLBase.Host = tenmanga.Hosts[0]

	RegisterService("tenmanga", tenmanga)
}

type TenMangaService Service

func (service *TenMangaService) Supports(u *url.URL) bool {
	return utils.StringSliceContains(tenmanga.Hosts, u.Host)
}

func (service *TenMangaService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	if serviceTenMangaRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			URL:     u,
			Service: service,
		}
		return manga, nil
	}

	if serviceTenMangaRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			URL:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *TenMangaService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceTenMangaHTMLSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 2 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[1]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data

	return name, nil
}

func (service *TenMangaService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceTenMangaHTMLSelectorMangaChapters.Find(rootNode)
	chapterCount := len(linkNodes)
	chapters := make([]*Chapter, 0, chapterCount)
	for i := chapterCount - 1; i >= 0; i-- {
		linkNode := linkNodes[i]
		chapterURL := utils.URLCopy(tenmanga.URLBase)
		chapterURL.Path = utils.HTMLGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			URL:     chapterURL,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *TenMangaService) ChapterName(chapter *Chapter) (string, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	mangaNameNodes := serviceTenMangaHTMLSelectorChapterNameManga.Find(rootNode)
	if len(mangaNameNodes) != 2 {
		return "", errors.New("Manga name node not found")
	}
	mangaNameNode := mangaNameNodes[1]
	mangaName, err := utils.HTMLGetNodeText(mangaNameNode)
	if err != nil {
		return "", err
	}

	titleNodes := serviceTenMangaHTMLSelectorChapterNameTitle.Find(rootNode)
	if len(titleNodes) != 1 {
		return "", errors.New("Title node not found")
	}
	titleNode := titleNodes[0]
	title, err := utils.HTMLGetNodeText(titleNode)
	if err != nil {
		return "", err
	}

	matches := serviceTenMangaRegexpChapterName.FindStringSubmatch(title)
	if matches == nil {
		return "", errors.New("Invalid title format")
	}
	name := matches[1]
	name = strings.TrimPrefix(name, mangaName)
	name = strings.TrimSpace(name)

	return name, nil
}

func (service *TenMangaService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	pageNodes := serviceTenMangaHTMLSelectorChapterPages.Find(rootNode)
	pages := make([]*Page, 0, len(pageNodes))
	for _, pageNode := range pageNodes {
		pageURL, err := url.Parse(utils.HTMLGetNodeAttribute(pageNode, "value"))
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

func (service *TenMangaService) PageImageURL(page *Page) (*url.URL, error) {
	rootNode, err := utils.HTTPGetHTML(page.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceTenMangaHTMLSelectorPageImage.Find(rootNode)
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

func (service *TenMangaService) HTTPRetry() int {
	return service.httpRetry
}
func (service *TenMangaService) SetHTTPRetry(nr int) {
	service.httpRetry = nr
}

func (service *TenMangaService) String() string {
	return "TenMangaService"
}
