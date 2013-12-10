package service

import (
	"net/url"
)

type Page struct {
	URL     *url.URL
	Service ServiceHandler
}

func (page *Page) ImageURL() (*url.URL, error) {
	return page.Service.PageImageURL(page)
}
