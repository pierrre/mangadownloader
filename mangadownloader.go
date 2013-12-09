package mangadownloader

import (
	"github.com/matrixik/mangadownloader/service"

	"fmt"
	"github.com/pierrre/archivefile/zip"
	"io/ioutil"
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

func init() {
	filenameCleanReplacements := make([]string, len(filenameReservedCharacters)*2)
	for _, char := range filenameReservedCharacters {
		filenameCleanReplacements = append(filenameCleanReplacements, string(char))
		filenameCleanReplacements = append(filenameCleanReplacements, " ")
	}
	filenameCleanReplacer = strings.NewReplacer(filenameCleanReplacements...)
}

type MangaDownloader struct {
	Services service.Services
}

func NewMangaDownloader() *MangaDownloader {
	md := new(MangaDownloader)
	md.Services = make(service.Services)

	return md
}

func CreateDefaultMangeDownloader() *MangaDownloader {
	md := new(MangaDownloader)

	//for name, _ := range md.Services {
	//	md.Services[name].{
	//	HttpRetry: 5,
	//	}
	//}

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

func (md *MangaDownloader) DownloadManga(manga *service.Manga, out string, options *Options) error {
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

func (md *MangaDownloader) downloadChapters(chapters []*service.Chapter, out string, options *Options) error {
	work := make(chan *service.Chapter)
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

	errs := make(service.MultiError, 0)
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

func (md *MangaDownloader) DownloadChapter(chapter *service.Chapter, out string, options *Options) error {
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

func (md *MangaDownloader) downloadChapter(chapter *service.Chapter, out string, options *Options) error {
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

func (md *MangaDownloader) downloadChapterCbz(chapter *service.Chapter, out string, options *Options) error {
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

func (md *MangaDownloader) downloadPages(pages []*service.Page, out string, options *Options) error {
	type pageWork struct {
		page  *service.Page
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

	errs := make(service.MultiError, 0)
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

func (md *MangaDownloader) downloadPageWithIndex(page *service.Page, out string, index int, options *Options) error {
	filenameFormat := "%0" + strconv.Itoa(options.PageDigitCount) + "d"
	filename := fmt.Sprintf(filenameFormat, index+1)
	return md.DownloadPage(page, out, filename, options)
}

func (md *MangaDownloader) DownloadPage(page *service.Page, out string, filename string, options *Options) error {
	out = filepath.Join(out, filename)

	imageUrl, err := page.ImageUrl()
	if err != nil {
		return err
	}

	response, err := service.HttpGet(imageUrl, 5)
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

type Options struct {
	Cbz             bool
	PageDigitCount  int
	ParallelChapter int
	ParallelPage    int
}

func cleanFilename(name string) string {
	return filenameCleanReplacer.Replace(name)
}
