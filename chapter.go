package mangadownloader

import (
	"net/url"
)

type Chapter struct {
	Url     *url.URL
	Service Service
}

func (chapter *Chapter) Name() (string, error) {
	return chapter.Service.ChapterName(chapter)
}

func (chapter *Chapter) Pages() ([]*Page, error) {
	return chapter.Service.ChapterPages(chapter)
}
