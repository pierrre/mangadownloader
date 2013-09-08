package mangadownloader

import (
	"net/url"
)

type Page struct {
	Url     *url.URL
	Service Service
}

func (page *Page) ImageUrl() (*url.URL, error) {
	return page.Service.PageImageUrl(page)
}
