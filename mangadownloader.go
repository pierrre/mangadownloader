package mangadownloader

import (
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var (
	regexpImageContentType *regexp.Regexp
)

func init() {
	regexpImageContentType, _ = regexp.Compile("^image/(.+)$")
}

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

func (md *MangaDownloader) DownloadManga(manga *Manga, out string) error {
	name, err := manga.Name()
	if err != nil {
		return err
	}
	out = filepath.Join(out, name)
	chapters, err := manga.Chapters()
	if err != nil {
		return err
	}
	for _, chapter := range chapters {
		err := md.DownloadChapter(chapter, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (md *MangaDownloader) DownloadChapter(chapter *Chapter, out string) error {
	name, err := chapter.Name()
	if err != nil {
		return err
	}
	out = filepath.Join(out, name)
	pages, err := chapter.Pages()
	if err != nil {
		return err
	}
	for _, page := range pages {
		err := md.DownloadPage(page, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (md *MangaDownloader) DownloadPage(page *Page, out string) error {
	index, err := page.Index()
	if err != nil {
		return err
	}
	out = filepath.Join(out, fmt.Sprintf("%04d", index))

	image, err := page.Image()
	if err != nil {
		return err
	}
	var response *http.Response
	for {
		response, err = md.HttpGet(image.Url)
		if err == nil {
			break
		}
		fmt.Println(err)
	}
	defer response.Body.Close()
	var extension string
	if len(extension) == 0 {
		contentType := response.Header.Get("content-type")
		if len(contentType) > 0 {
			matches := regexpImageContentType.FindStringSubmatch(contentType)
			if matches != nil && len(matches) == 2 {
				extension = matches[1]
			}
		}
	}
	if len(extension) > 0 {
		if extension == "jpeg" {
			extension = "jpg"
		}
		out += "." + extension
	}
	err = os.MkdirAll(filepath.Dir(out), 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
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
