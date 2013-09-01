package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
)

const (
	serviceMangaReaderDomain     = "www.mangareader.net"
	serviceMangaReaderPathMangas = "/alphabetical"
)

var (
	serviceMangaReaderUrlBase   *url.URL
	serviceMangaReaderUrlMangas *url.URL

	serviceMangaReaderHtmlSelectorMangas   *selector.Chain
	serviceMangaReaderHtmlSelectorChapters *selector.Chain
	serviceMangaReaderHtmlSelectorPages    *selector.Chain
	serviceMangaReaderHtmlSelectorImage    *selector.Chain
)

func init() {
	serviceMangaReaderUrlBase = new(url.URL)
	serviceMangaReaderUrlBase.Scheme = "http"
	serviceMangaReaderUrlBase.Host = serviceMangaReaderDomain

	serviceMangaReaderUrlMangas = urlCopy(serviceMangaReaderUrlBase)
	serviceMangaReaderUrlMangas.Path = serviceMangaReaderPathMangas

	serviceMangaReaderHtmlSelectorMangas, _ = selector.Selector("ul.series_alpha a")

	serviceMangaReaderHtmlSelectorChapters, _ = selector.Selector("#listing a")

	serviceMangaReaderHtmlSelectorPages, _ = selector.Selector("#pageMenu option")

	serviceMangaReaderHtmlSelectorImage, _ = selector.Selector("#img")
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

func (service *MangaReaderService) Chapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := service.Md.HttpGetHtml(manga.Url)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaReaderHtmlSelectorChapters.Find(rootNode)

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

func (service *MangaReaderService) Pages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := service.Md.HttpGetHtml(chapter.Url)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaReaderHtmlSelectorPages.Find(rootNode)

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

func (service *MangaReaderService) Image(page *Page) (*Image, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaReaderHtmlSelectorImage.Find(rootNode)
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
