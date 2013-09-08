package mangadownloader

import (
	"net/url"
)

type Page struct {
	Url     *url.URL
	Service Service
}

func (page *Page) Image() (*Image, error) {
	return page.Service.PageImage(page)
}
