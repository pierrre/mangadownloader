package mangadownloader

import (
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

	"github.com/pierrre/archivefile/zip"
	"golang.org/x/net/html"
)

var (
	regexpImageContentType, _ = regexp.Compile("^image/(.+)$")
	filenameCleanReplacer     *strings.Replacer
)

func init() {
	filenameCleanReplacements := make([]string, len(filenameReservedCharacters)*2)
	for _, char := range filenameReservedCharacters {
		filenameCleanReplacements = append(filenameCleanReplacements, string(char))
		filenameCleanReplacements = append(filenameCleanReplacements, " ")
	}
	filenameCleanReplacer = strings.NewReplacer(filenameCleanReplacements...)
}

type MangaDownloader struct {
	Services  map[string]Service
	HttpRetry int
}

func NewMangaDownloader() *MangaDownloader {
	md := new(MangaDownloader)
	md.Services = make(map[string]Service)

	return md
}

func CreateDefaultMangeDownloader() *MangaDownloader {
	md := NewMangaDownloader()

	md.HttpRetry = 5

	md.Services["mangafox"] = &MangaFoxService{
		Md: md,
	}

	md.Services["mangahere"] = &MangaHereService{
		Md: md,
	}

	md.Services["mangareader"] = &MangaReaderService{
		Md: md,
	}

	md.Services["mangawall"] = &MangaWallService{
		Md: md,
	}

	md.Services["tenmanga"] = &TenMangaService{
		Md: md,
	}

	return md
}

func (md *MangaDownloader) Identify(u *url.URL) (interface{}, error) {
	for _, service := range md.Services {
		if service.Supports(u) {
			return service.Identify(u)
		}
	}

	return nil, fmt.Errorf("url '%s' not supported by any service", u)
}

func (md *MangaDownloader) DownloadManga(manga *Manga, out string, options *Options) error {
	name, err := manga.Name()
	if err != nil {
		return err
	}

	out = filepath.Join(out, cleanFilename(name))

	chapters, err := manga.Chapters()
	if err != nil {
		return err
	}

	err = md.downloadChapters(chapters, out, options)
	if err != nil {
		return err
	}

	return nil
}

func (md *MangaDownloader) downloadChapters(chapters []*Chapter, out string, options *Options) error {
	work := make(chan *Chapter)
	go func() {
		for _, chapter := range chapters {
			work <- chapter
		}
		close(work)
	}()

	parallelChapter := options.ParallelChapter
	if parallelChapter < 1 {
		parallelChapter = 1
	}
	wg := new(sync.WaitGroup)
	wg.Add(parallelChapter)
	result := make(chan error)
	for i := 0; i < parallelChapter; i++ {
		go func() {
			for chapter := range work {
				result <- md.DownloadChapter(chapter, out, options)
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

func (md *MangaDownloader) DownloadChapter(chapter *Chapter, out string, options *Options) error {
	name, err := chapter.Name()
	if err != nil {
		return err
	}
	out = filepath.Join(out, cleanFilename(name))

	if options.Cbz {
		return md.downloadChapterCbz(chapter, out, options)
	} else {
		return md.downloadChapter(chapter, out, options)
	}
}

func (md *MangaDownloader) downloadChapter(chapter *Chapter, out string, options *Options) error {
	if fileExists(out) {
		return nil
	}

	outTmp := out + ".tmp"
	if fileExists(outTmp) {
		err := os.RemoveAll(outTmp)
		if err != nil {
			return err
		}
	}

	pages, err := chapter.Pages()
	if err != nil {
		return err
	}

	err = md.downloadPages(pages, outTmp, options)
	if err != nil {
		return err
	}

	err = os.Rename(outTmp, out)
	if err != nil {
		return err
	}

	return nil
}

func (md *MangaDownloader) downloadChapterCbz(chapter *Chapter, out string, options *Options) error {
	outCbz := out + ".cbz"
	if fileExists(outCbz) {
		return nil
	}

	outCbzTmp := outCbz + ".tmp"
	if fileExists(outCbzTmp) {
		err := os.RemoveAll(outCbzTmp)
		if err != nil {
			return err
		}
	}

	err := md.downloadChapter(chapter, out, options)
	if err != nil {
		return err
	}

	err = zip.ArchiveFile(out+string(filepath.Separator), outCbzTmp, nil)
	if err != nil {
		return err
	}

	err = os.Rename(outCbzTmp, outCbz)
	if err != nil {
		return err
	}

	err = os.RemoveAll(out)
	if err != nil {
		return err
	}

	return nil
}

func (md *MangaDownloader) downloadPages(pages []*Page, out string, options *Options) error {
	type pageWork struct {
		page  *Page
		index int
	}

	work := make(chan *pageWork)
	go func() {
		for index, page := range pages {
			work <- &pageWork{
				page:  page,
				index: index,
			}
		}
		close(work)
	}()

	parallelPage := options.ParallelPage
	if parallelPage < 1 {
		parallelPage = 1
	}
	wg := new(sync.WaitGroup)
	wg.Add(parallelPage)
	result := make(chan error)
	for i := 0; i < parallelPage; i++ {
		go func() {
			for chapterPageWork := range work {
				result <- md.downloadPageWithIndex(chapterPageWork.page, out, chapterPageWork.index, options)
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

func (md *MangaDownloader) downloadPageWithIndex(page *Page, out string, index int, options *Options) error {
	filenameFormat := "%0" + strconv.Itoa(options.PageDigitCount) + "d"
	filename := fmt.Sprintf(filenameFormat, index+1)
	return md.DownloadPage(page, out, filename, options)
}

func (md *MangaDownloader) DownloadPage(page *Page, out string, filename string, options *Options) error {
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
			if matches != nil {
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
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.4 Safari/537.36")

	httpRetry := md.HttpRetry
	if httpRetry < 1 {
		httpRetry = 1
	}

	errs := make(MultiError, 0)
	for i := 0; i < httpRetry; i++ {
		response, err := http.DefaultClient.Do(request)
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

type Options struct {
	Cbz             bool
	PageDigitCount  int
	ParallelChapter int
	ParallelPage    int
}

func cleanFilename(name string) string {
	return filenameCleanReplacer.Replace(name)
}
