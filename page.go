package mangadownloader

import (
	"net/url"
)

type Page struct {
	Url     *url.URL
	Service Service
}
