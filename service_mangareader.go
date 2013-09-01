package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"net/url"
)

type MangaReaderService struct {
	Md *MangaDownloader
}

func (service *MangaReaderService) Mangas() ([]*Manga, error) {
	u, err := url.Parse("http://www.mangareader.net/alphabetical")
	if err != nil {
		return nil, err
	}
	rootNode, err := service.Md.HttpGetHtml(u)
	if err != nil {
		return nil, err
	}

	linkSelector, err := selector.Selector("ul.series_alpha a")
	if err != nil {
		return nil, err
	}
	linkNodes := linkSelector.Find(rootNode)

	mangas := make([]*Manga, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		href := getHtmlNodeAttribute(linkNode, "href")
		mangaUrl, err := url.Parse("http://www.mangareader.net" + href)
		if err != nil {
			return nil, err
		}
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

	linkSelector, err := selector.Selector("#listing a")
	if err != nil {
		return nil, err
	}
	linkNodes := linkSelector.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		href := getHtmlNodeAttribute(linkNode, "href")
		chapterUrl, err := url.Parse("http://www.mangareader.net" + href)
		if err != nil {
			return nil, err
		}
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

	optionSelector, err := selector.Selector("#pageMenu option")
	if err != nil {
		return nil, err
	}
	optionNodes := optionSelector.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		value := getHtmlNodeAttribute(optionNode, "value")
		pageUrl, err := url.Parse("http://www.mangareader.net" + value)
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

func (service *MangaReaderService) Image(page *Page) (*Image, error) {
	rootNode, err := service.Md.HttpGetHtml(page.Url)
	if err != nil {
		return nil, err
	}

	imgSelector, err := selector.Selector("#img")
	if err != nil {
		return nil, err
	}
	imgNodes := imgSelector.Find(rootNode)
	if len(imgNodes) < 1 {
		return nil, errors.New("Image node not found")
	}
	imgNode := imgNodes[0]

	var src string
	for _, attr := range imgNode.Attr {
		if attr.Key == "src" {
			src = attr.Val
		}
	}
	imageUrl, err := url.Parse(src)
	image := &Image{
		Url:     imageUrl,
		Service: service,
	}

	return image, nil
}

func (service *MangaReaderService) String() string {
	return "MangaReaderService"
}
