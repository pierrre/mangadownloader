package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
	"strconv"
)

const (
	serviceMangaReaderDomain     = "www.mangareader.net"
	serviceMangaReaderPathMangas = "/alphabetical"
)

var (
	serviceMangaReaderUrlBase   *url.URL
	serviceMangaReaderUrlMangas *url.URL

	serviceMangaReaderHtmlSelectorMangas        *selector.Chain
	serviceMangaReaderHtmlSelectorMangaName     *selector.Chain
	serviceMangaReaderHtmlSelectorMangaChapters *selector.Chain
	serviceMangaReaderHtmlSelectorChapterPages  *selector.Chain
	serviceMangaReaderHtmlSelectorPageIndex     *selector.Chain
	serviceMangaReaderHtmlSelectorPageImage     *selector.Chain
)

func init() {
	serviceMangaReaderUrlBase = new(url.URL)
	serviceMangaReaderUrlBase.Scheme = "http"
	serviceMangaReaderUrlBase.Host = serviceMangaReaderDomain

	serviceMangaReaderUrlMangas = urlCopy(serviceMangaReaderUrlBase)
	serviceMangaReaderUrlMangas.Path = serviceMangaReaderPathMangas

	serviceMangaReaderHtmlSelectorMangas, _ = selector.Selector("ul.series_alpha a")

	serviceMangaReaderHtmlSelectorMangaName, _ = selector.Selector("h2.aname")

	serviceMangaReaderHtmlSelectorMangaChapters, _ = selector.Selector("#chapterlist a")

	serviceMangaReaderHtmlSelectorChapterPages, _ = selector.Selector("#pageMenu option")

	serviceMangaReaderHtmlSelectorPageIndex, _ = selector.Selector("select#pageMenu option[selected=selected]")

	serviceMangaReaderHtmlSelectorPageImage, _ = selector.Selector("#img")
}

type MangaReaderService struct {
	Md *MangaDownloader
}

func (service *MangaReaderService) Mangas() ([]*Manga, error) {
	rootNode, err := service.Md.HttpGetHtml(serviceMangaReaderUrlMangas)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaReaderHtmlSelectorMangas.Find(rootNode)

	mangas := make([]*Manga, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		mangaUrl := urlCopy(serviceMangaReaderUrlBase)
		mangaUrl.Path = htmlGetNodeAttribute(linkNode, "href")
		manga := &Manga{
			Url:     mangaUrl,
			Service: service,
		}
		mangas = append(mangas, manga)
	}

	return mangas, nil
}

func (service *MangaReaderService) MangaName(manga *Manga) (string, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaReaderHtmlSelectorMangaName.Find(rootNode)
	if len(nameNodes) < 1 {
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
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaReaderHtmlSelectorMangaChapters.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl := urlCopy(serviceMangaReaderUrlBase)
		chapterUrl.Path = htmlGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *MangaReaderService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaReaderHtmlSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageUrl := urlCopy(serviceMangaReaderUrlBase)
		pageUrl.Path = htmlGetNodeAttribute(optionNode, "value")
		page := &Page{
			Url:     pageUrl,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaReaderService) PageIndex(page *Page) (uint, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return 0, err
	}

	indexNodes := serviceMangaReaderHtmlSelectorPageIndex.Find(rootNode)
	if len(indexNodes) < 1 {
		return 0, errors.New("Index node not found")
	}
	indexNode := indexNodes[0]
	if indexNode.FirstChild == nil {
		return 0, errors.New("Index text node not found")
	}
	indexTextNode := indexNode.FirstChild
	indexString := indexTextNode.Data
	indexInt, err := strconv.Atoi(indexString)
	if err != nil {
		return 0, err
	}
	index := uint(indexInt)

	return index, err
}

func (service *MangaReaderService) PageImage(page *Page) (*Image, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaReaderHtmlSelectorPageImage.Find(rootNode)
	if len(imgNodes) < 1 {
		return nil, errors.New("Image node not found")
	}
	imgNode := imgNodes[0]

	imageUrl, err := url.Parse(htmlGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}
	image := &Image{
		Url:     imageUrl,
		Service: service,
	}

	return image, nil
}

func (service *MangaReaderService) String() string {
	return "MangaReaderService"
}
