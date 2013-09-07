package mangadownloader

import (
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

const (
	httpRetry = 5
)

var (
	regexpImageContentType, _  = regexp.Compile("^image/(.+)$")
	filenameReservedCharacters = []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}
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

	work := make(chan *Chapter)
	go func() {
		for _, chapter := range chapters {
			work <- chapter
		}
		close(work)
	}()

	concurrent := 4
	wg := new(sync.WaitGroup)
	wg.Add(concurrent)
	result := make(chan error)
	for i := 0; i < concurrent; i++ {
		go func() {
			for chapter := range work {
				result <- md.DownloadChapter(chapter, out)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	errs := make(MultiError, 0)
	for err := range result {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (md *MangaDownloader) DownloadChapter(chapter *Chapter, out string) error {
	name, err := chapter.Name()
	if err != nil {
		return err
	}
	out = filepath.Join(out, name)
	if fileExists(out) {
		return nil
	}
	outTmp := out + ".tmp"
	if fileExists(outTmp) {
		err = os.RemoveAll(outTmp)
		if err != nil {
			return err
		}
	}

	pages, err := chapter.Pages()
	if err != nil {
		return err
	}

	work := make(chan *Page)
	go func() {
		for _, page := range pages {
			work <- page
		}
		close(work)
	}()

	concurrent := 8
	wg := new(sync.WaitGroup)
	wg.Add(concurrent)
	result := make(chan error)
	for i := 0; i < concurrent; i++ {
		go func() {
			for page := range work {
				result <- md.DownloadPage(page, outTmp)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	errs := make(MultiError, 0)
	for err := range result {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	err = os.Rename(outTmp, out)
	if err != nil {
		return err
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

	response, err := md.HttpGet(image.Url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

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
	err = ioutil.WriteFile(out, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (md *MangaDownloader) HttpGet(u *url.URL) (response *http.Response, err error) {
	return md.httpGetRetry(u, httpRetry)
}

func (md *MangaDownloader) httpGetRetry(u *url.URL, retry int) (*http.Response, error) {
	if retry < 1 {
		return nil, errors.New("Invalid retry")
	}

	errs := make(MultiError, 0)
	for i := 0; i < retry; i++ {
		response, err := http.Get(u.String())
		if err == nil {
			return response, nil
		}
		errs = append(errs, err)
	}
	return nil, errs
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
