package mangadownloader

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"net/http"
	"net/url"
)

type Service struct {
	Url *url.URL
}

func (service *Service) Mangas() ([]*Manga, error) {
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
			Url: mangaUrl,
		}
		mangas = append(mangas, manga)
	}

	return mangas, nil
}
