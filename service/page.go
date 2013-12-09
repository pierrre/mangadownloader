package service

import (
	"net/url"
)

type Page struct {
	Url     *url.URL
	Service ServiceHandler
}

func (page *Page) ImageUrl() (*url.URL, error) {
	return page.Service.PageImageUrl(page)
}
