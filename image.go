package mangadownloader

import (
	"net/url"
)

type Image struct {
	Url     *url.URL
	Service Service
}
