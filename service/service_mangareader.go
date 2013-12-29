package service

import (
	"github.com/matrixik/mangadownloader/utils"

	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

const (
	serviceMangaReadeChapterDigitCount = 4
)

var (
	mangareader = &MangaReaderService{}

	serviceMangaReaderHTMLSelectorIdentifyManga, _   = selector.Selector("#chapterlist")
	serviceMangaReaderHTMLSelectorIdentifyChapter, _ = selector.Selector("#pageMenu")
	serviceMangaReaderHTMLSelectorMangaName, _       = selector.Selector("h2.aname")
	serviceMangaReaderHTMLSelectorMangaChapters, _   = selector.Selector("#chapterlist a")
	serviceMangaReaderHTMLSelectorChapterName, _     = selector.Selector("#mangainfo h1")
	serviceMangaReaderHTMLSelectorChapterPages, _    = selector.Selector("#pageMenu option")
	serviceMangaReaderHTMLSelectorPageImage, _       = selector.Selector("#img")

	serviceMangaReaderRegexpChapterName, _ = regexp.Compile("^.* (\\d*)$")

	serviceMangaReaderFormatChapter = "%0" + strconv.Itoa(serviceMangaReadeChapterDigitCount) + "d"
)

func init() {
	mangareader.ServiceCommon = ServiceCommon{
		Hosts: []string{
			"www.mangareader.net",
			"mangareader.net",
		},
	}

	mangareader.URLBase = new(url.URL)
	mangareader.URLBase.Scheme = "http"
	mangareader.URLBase.Host = mangareader.Hosts[0]

	RegisterService("mangareader", mangareader)
}

type MangaReaderService struct {
	ServiceCommon
}

func init() {
	RegisterService("mangareader", mangareader)
}

func (service *MangaReaderService) Supports(u *url.URL) bool {
	return utils.StringSliceContains(mangareader.Hosts, u.Host)
}

func (service *MangaReaderService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	rootNode, err := utils.HTTPGetHTML(u, service.httpRetry)
	if err != nil {
		return nil, err
	}

	identifyMangaNodes := serviceMangaReaderHTMLSelectorIdentifyManga.Find(rootNode)
	if len(identifyMangaNodes) == 1 {
		manga := &Manga{
			URL:     u,
			Service: service,
		}
		return manga, nil
	}

	identifyChapterNodes := serviceMangaReaderHTMLSelectorIdentifyChapter.Find(rootNode)
	if len(identifyChapterNodes) == 1 {
		chapter := &Chapter{
			URL:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *MangaReaderService) MangaName(manga *Manga) (string, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaReaderHTMLSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data

	return name, nil
}

func (service *MangaReaderService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := utils.HTTPGetHTML(manga.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaReaderHTMLSelectorMangaChapters.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterURL := utils.URLCopy(mangareader.URLBase)
		chapterURL.Path = utils.HTMLGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			URL:     chapterURL,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *MangaReaderService) ChapterName(chapter *Chapter) (string, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaReaderHTMLSelectorChapterName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data
	matches := serviceMangaReaderRegexpChapterName.FindStringSubmatch(name)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name = matches[1]
	nameInt, err := strconv.Atoi(name)
	if err != nil {
		return "", err
	}
	name = fmt.Sprintf(serviceMangaReaderFormatChapter, nameInt)

	return name, nil
}

func (service *MangaReaderService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := utils.HTTPGetHTML(chapter.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaReaderHTMLSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageURL := utils.URLCopy(mangareader.URLBase)
		pageURL.Path = utils.HTMLGetNodeAttribute(optionNode, "value")
		page := &Page{
			URL:     pageURL,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaReaderService) PageImageURL(page *Page) (*url.URL, error) {
	rootNode, err := utils.HTTPGetHTML(page.URL, service.httpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaReaderHTMLSelectorPageImage.Find(rootNode)
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

func (service *MangaReaderService) HTTPRetry() int {
	return service.httpRetry
}
func (service *MangaReaderService) SetHTTPRetry(nr int) {
	service.httpRetry = nr
}

func (service *MangaReaderService) String() string {
	return "MangaReaderService"
}
