package service

import (
	"net/url"
)

type Chapter struct {
	Url     *url.URL
	Service ServiceHandler
}

func (chapter *Chapter) Name() (string, error) {
	return chapter.Service.ChapterName(chapter)
}

func (chapter *Chapter) Pages() ([]*Page, error) {
	return chapter.Service.ChapterPages(chapter)
}
