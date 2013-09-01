package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"errors"
	"net/http"
	"net/url"
)

type MangaReaderService struct {
}

func (service *MangaReaderService) Mangas() ([]*Manga, error) {
	response, err := http.Get("http://www.mangareader.net/alphabetical")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rootNode, err := html.Parse(response.Body)
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
		var href string
		for _, attr := range linkNode.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if len(href) == 0 {
			continue
		}
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
	response, err := http.Get(manga.Url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rootNode, err := html.Parse(response.Body)
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
		var href string
		for _, attr := range linkNode.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		if len(href) == 0 {
			continue
		}
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
	response, err := http.Get(chapter.Url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rootNode, err := html.Parse(response.Body)
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
		var value string
		for _, attr := range optionNode.Attr {
			if attr.Key == "value" {
				value = attr.Val
			}
		}
		if len(value) == 0 {
			continue
		}
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
	response, err := http.Get(page.Url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rootNode, err := html.Parse(response.Body)
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
