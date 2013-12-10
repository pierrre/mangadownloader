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

func chapterSliceReverse(chapters []*Chapter) []*Chapter {
	count := len(chapters)
	reversed := make([]*Chapter, 0, count)
	for i := count - 1; i >= 0; i-- {
		reversed = append(reversed, chapters[i])
	}
	return reversed
}
