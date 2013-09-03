package mangadownloader

import (
	"net/url"
)

type Chapter struct {
	Url     *url.URL
	Service Service
}

func (chapter *Chapter) Pages() ([]*Page, error) {
	return chapter.Service.ChapterPages(chapter)
}
