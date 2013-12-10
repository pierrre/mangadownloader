package service

import (
	"net/url"
)

type Service struct {
	Hosts     []string
	URLBase   *url.URL
	httpRetry int
}

type ServiceHandler interface {
	Supports(*url.URL) bool
	Identify(*url.URL) (interface{}, error)
	MangaName(*Manga) (string, error)
	MangaChapters(*Manga) ([]*Chapter, error)
	ChapterName(*Chapter) (string, error)
	ChapterPages(*Chapter) ([]*Page, error)
	PageImageURL(*Page) (*url.URL, error)

	HTTPRetry() int
	SetHTTPRetry(int)
}

var Services = map[string]ServiceHandler{}

func RegisterService(name string, service ServiceHandler) {
	Services[name] = service
}

func FindService(name string) (c ServiceHandler, ok bool) {
	c, ok = Services[name]
	return
}
