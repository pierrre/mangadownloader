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
	"strconv"
	"strings"
	"sync"
)

var (
	regexpImageContentType, _ = regexp.Compile("^image/(.+)$")
	filenameCleanReplacer     *strings.Replacer
)

type chapterPageWork struct {
	page  *Page
	index uint
}

func init() {
	filenameCleanReplacements := make([]string, len(filenameReservedCharacters)*2)
	for _, char := range filenameReservedCharacters {
		filenameCleanReplacements = append(filenameCleanReplacements, string(char))
		filenameCleanReplacements = append(filenameCleanReplacements, " ")
	}
	filenameCleanReplacer = strings.NewReplacer(filenameCleanReplacements...)
}

type MangaDownloader struct {
	Services        []Service
	PageDigitCount  int
	HttpRetry       int
	ParallelChapter int
	ParallelPage    int
}

func CreateDefaultMangeDownloader() *MangaDownloader {
	md := &MangaDownloader{}

	md.Services = append(md.Services, &MangaReaderService{
		Md: md,
	})

	md.Services = append(md.Services, &MangaFoxService{
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
	out = filepath.Join(out, md.cleanFilename(name))
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

	parallelChapter := md.ParallelChapter
	if parallelChapter < 1 {
		parallelChapter = 1
	}
	wg := new(sync.WaitGroup)
	wg.Add(parallelChapter)
	result := make(chan error)
	for i := 0; i < parallelChapter; i++ {
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
	out = filepath.Join(out, md.cleanFilename(name))
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

	work := make(chan *chapterPageWork)
	go func() {
		for index, page := range pages {
			work <- &chapterPageWork{
				page:  page,
				index: uint(index),
			}
		}
		close(work)
	}()

	parallelPage := md.ParallelPage
	if parallelPage < 1 {
		parallelPage = 1
	}
	wg := new(sync.WaitGroup)
	wg.Add(parallelPage)
	result := make(chan error)
	for i := 0; i < parallelPage; i++ {
		go func() {
			for chapterPageWork := range work {
				result <- md.downloadPageWithIndex(chapterPageWork.page, outTmp, chapterPageWork.index)
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

func (md *MangaDownloader) downloadPageWithIndex(page *Page, out string, index uint) error {
	filenameFormat := "%0" + strconv.Itoa(md.PageDigitCount) + "d"
	filename := fmt.Sprintf(filenameFormat, index+1)
	return md.DownloadPage(page, out, filename)
}

func (md *MangaDownloader) DownloadPage(page *Page, out string, filename string) error {
	out = filepath.Join(out, filename)

	imageUrl, err := page.ImageUrl()
	if err != nil {
		return err
	}

	response, err := md.HttpGet(imageUrl)
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
	httpRetry := md.HttpRetry
	if httpRetry < 1 {
		httpRetry = 1
	}

	errs := make(MultiError, 0)
	for i := 0; i < httpRetry; i++ {
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

func (md *MangaDownloader) cleanFilename(file string) string {
	return filenameCleanReplacer.Replace(file)
}
