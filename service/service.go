package service

import (
	"net/url"
)

type ServiceCommon struct {
	Hosts     []string
	URLBase   *url.URL
	httpRetry int
}

type Service interface {
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

var Services = map[string]Service{}

func RegisterService(name string, service Service) {
	Services[name] = service
}

func FindService(name string) (c Service, ok bool) {
	c, ok = Services[name]
	return
}
