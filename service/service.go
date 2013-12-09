package service

import (
	"net/url"
)

type Service struct {
	Hosts     []string
	UrlBase   *url.URL
	HttpRetry int
}

type ServiceHandler interface {
	Supports(*url.URL) bool
	Identify(*url.URL) (interface{}, error)
	MangaName(*Manga) (string, error)
	MangaChapters(*Manga) ([]*Chapter, error)
	ChapterName(*Chapter) (string, error)
	ChapterPages(*Chapter) ([]*Page, error)
	PageImageUrl(*Page) (*url.URL, error)
}

type Services map[string]ServiceHandler

var services Services

func RegisterService(name string, service ServiceHandler) {
	services[name] = service
}

func FindService(name string) (c ServiceHandler, ok bool) {
	c, ok = services[name]
	return
}
