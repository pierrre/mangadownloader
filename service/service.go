package service

import (
	"net/url"
)

type Service struct {
	Hosts     []string
	UrlBase   *url.URL
	httpRetry int
}

type ServiceHandler interface {
	Supports(*url.URL) bool
	Identify(*url.URL) (interface{}, error)
	MangaName(*Manga) (string, error)
	MangaChapters(*Manga) ([]*Chapter, error)
	ChapterName(*Chapter) (string, error)
	ChapterPages(*Chapter) ([]*Page, error)
	PageImageUrl(*Page) (*url.URL, error)

	HttpRetry() int
	SetHttpRetry(int)
}

var Services = map[string]ServiceHandler{}

func RegisterService(name string, service ServiceHandler) {
	Services[name] = service
}

func FindService(name string) (c ServiceHandler, ok bool) {
	c, ok = Services[name]
	return
}
