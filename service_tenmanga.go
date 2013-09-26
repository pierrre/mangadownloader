package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

const (
	serviceTenMangaDomain = "www.tenmanga.com"
)

var (
	serviceTenMangaUrlBase *url.URL

	serviceTenMangaHtmlSelectorMangaName, _        = selector.Selector(".postion .red")
	serviceTenMangaHtmlSelectorMangaChapters, _    = selector.Selector(".chapter_list td[align=left] a")
	serviceTenMangaHtmlSelectorChapterNameTitle, _ = selector.Selector("title")
	serviceTenMangaHtmlSelectorChapterNameManga, _ = selector.Selector(".postion a")
	serviceTenMangaHtmlSelectorChapterPages, _     = selector.Selector("#page option")
	serviceTenMangaHtmlSelectorPageImage, _        = selector.Selector("#comicpic")

	serviceTenMangaRegexpIdentifyManga, _   = regexp.Compile("^/book/.+$")
	serviceTenMangaRegexpIdentifyChapter, _ = regexp.Compile("^/chapter/.+$")
	serviceTenMangaRegexpChapterName, _     = regexp.Compile("^(.+) Page \\d+$")
)

func init() {
	serviceTenMangaUrlBase = new(url.URL)
	serviceTenMangaUrlBase.Scheme = "http"
	serviceTenMangaUrlBase.Host = serviceTenMangaDomain
}

type TenMangaService struct {
	Md *MangaDownloader
}

func (service *TenMangaService) Supports(u *url.URL) bool {
	return u.Host == serviceTenMangaDomain
}

func (service *TenMangaService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	if serviceTenMangaRegexpIdentifyManga.MatchString(u.Path) {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	if serviceTenMangaRegexpIdentifyChapter.MatchString(u.Path) {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *TenMangaService) MangaName(manga *Manga) (string, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := serviceTenMangaHtmlSelectorMangaName.Find(rootNode)
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
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceTenMangaHtmlSelectorMangaChapters.Find(rootNode)
	chapterCount := len(linkNodes)
	chapters := make([]*Chapter, 0, chapterCount)
	for i := chapterCount - 1; i >= 0; i-- {
		linkNode := linkNodes[i]
		chapterUrl := urlCopy(serviceTenMangaUrlBase)
		chapterUrl.Path = htmlGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *TenMangaService) ChapterName(chapter *Chapter) (string, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return "", err
	}

	mangaNameNodes := serviceTenMangaHtmlSelectorChapterNameManga.Find(rootNode)
	if len(mangaNameNodes) != 2 {
		return "", errors.New("Manga name node not found")
	}
	mangaNameNode := mangaNameNodes[1]
	mangaName, err := htmlGetNodeText(mangaNameNode)
	if err != nil {
		return "", err
	}

	titleNodes := serviceTenMangaHtmlSelectorChapterNameTitle.Find(rootNode)
	if len(titleNodes) != 1 {
		return "", errors.New("Title node not found")
	}
	titleNode := titleNodes[0]
	title, err := htmlGetNodeText(titleNode)
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
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	pageNodes := serviceTenMangaHtmlSelectorChapterPages.Find(rootNode)
	pages := make([]*Page, 0, len(pageNodes))
	for _, pageNode := range pageNodes {
		pageUrl, err := url.Parse(htmlGetNodeAttribute(pageNode, "value"))
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

func (service *TenMangaService) PageImageUrl(page *Page) (*url.URL, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceTenMangaHtmlSelectorPageImage.Find(rootNode)
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

func (service *TenMangaService) String() string {
	return "TenMangaService"
}
