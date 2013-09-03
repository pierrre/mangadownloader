package mangadownloader

import (
	"net/url"
)

type Page struct {
	Url     *url.URL
	Service Service
}

func (page *Page) Index() (uint, error) {
	return page.Service.PageIndex(page)
}

func (page *Page) Image() (*Image, error) {
	return page.Service.Image(page)
}
