package mangadownloader

import (
	"code.google.com/p/go.net/html"
	"errors"
	"net/http"
	"net/url"
)

type MangaDownloader struct {
	Services []Service
}

func CreateDefaultMangeDownloader() *MangaDownloader {
	md := &MangaDownloader{}

	md.Services = append(md.Services, &MangaReaderService{
		Md: md,
	})

	return md
}

func (md *MangaDownloader) Identify(u *url.URL) (interface{}, error) {
	for _, service := range md.Services {
		if service.Supports(u) {
			return service.Identify(u)
		}
	}

	return nil, errors.New("Unsupported url")
}

func (md *MangaDownloader) HttpGet(u *url.URL) (*http.Response, error) {
	return http.Get(u.String())
}

func (md *MangaDownloader) HttpGetHtml(u *url.URL) (*html.Node, error) {
	response, err := md.HttpGet(u)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	node, err := html.Parse(response.Body)
	return node, err
}
