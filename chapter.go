package mangadownloader

import (
	"net/url"
)

type Chapter struct {
	Url     *url.URL
	Service Service
}
