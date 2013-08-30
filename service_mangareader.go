package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
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

func (service *MangaReaderService) String() string {
	return "MangaReaderService"
}
