package mangadownloader

import (
	"github.com/pierrre/mangadownloader/service"
	"github.com/pierrre/mangadownloader/utils"

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

type Options struct {
	Cbz             bool
	PageDigitCount  int
	ParallelChapter int
	ParallelPage    int

	HTTPRetry int
}

func init() {
	filenameCleanReplacements := make([]string, len(filenameReservedCharacters)*2)
	for _, char := range filenameReservedCharacters {
		filenameCleanReplacements = append(filenameCleanReplacements, string(char))
		filenameCleanReplacements = append(filenameCleanReplacements, " ")
	}
	filenameCleanReplacer = strings.NewReplacer(filenameCleanReplacements...)

	for _, serviceStruct := range service.Services {
		serviceStruct.SetHTTPRetry(5)
	}
}

func Identify(u *url.URL, options *Options) (interface{}, error) {
	for _, service := range service.Services {
		if service.Supports(u) {
			service.SetHTTPRetry(options.HTTPRetry)
			return service.Identify(u)
		}
	}

	return nil, fmt.Errorf("url '%s' not supported by any service", u)
}

func DownloadManga(manga *service.Manga, out string, options *Options) error {
	name, err := manga.Name()
	if err != nil {
		return err
	}

	out = filepath.Join(out, cleanFilename(name))

	chapters, err := manga.Chapters()
	if err != nil {
		return err
	}

	err = downloadChapters(chapters, out, options)
	if err != nil {
		return err
	}

	return nil
}

func downloadChapters(chapters []*service.Chapter, out string, options *Options) error {
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
				result <- DownloadChapter(chapter, out, options)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	errs := make(utils.MultiError, 0)
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

func DownloadChapter(chapter *service.Chapter, out string, options *Options) error {
	name, err := chapter.Name()
	if err != nil {
		return err
	}
	out = filepath.Join(out, cleanFilename(name))

	if options.Cbz {
		return downloadChapterCbz(chapter, out, options)
	}
	return downloadChapter(chapter, out, options)
}

func downloadChapter(chapter *service.Chapter, out string, options *Options) error {
	if utils.FileExists(out) {
		return nil
	}

	outTmp := out + ".tmp"
	if utils.FileExists(outTmp) {
		err := os.RemoveAll(outTmp)
		if err != nil {
			return err
		}
	}

	pages, err := chapter.Pages()
	if err != nil {
		return err
	}

	err = downloadPages(pages, outTmp, options)
	if err != nil {
		return err
	}

	err = os.Rename(outTmp, out)
	if err != nil {
		return err
	}

	return nil
}

func downloadChapterCbz(chapter *service.Chapter, out string, options *Options) error {
	outCbz := out + ".cbz"
	if utils.FileExists(outCbz) {
		return nil
	}

	outCbzTmp := outCbz + ".tmp"
	if utils.FileExists(outCbzTmp) {
		err := os.RemoveAll(outCbzTmp)
		if err != nil {
			return err
		}
	}

	err := downloadChapter(chapter, out, options)
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

func downloadPages(pages []*service.Page, out string, options *Options) error {
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
				result <- downloadPageWithIndex(chapterPageWork.page, out, chapterPageWork.index, options)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	errs := make(utils.MultiError, 0)
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

func downloadPageWithIndex(page *service.Page, out string, index int, options *Options) error {
	filenameFormat := "%0" + strconv.Itoa(options.PageDigitCount) + "d"
	filename := fmt.Sprintf(filenameFormat, index+1)
	return DownloadPage(page, out, filename, options)
}

func DownloadPage(page *service.Page, out string, filename string, options *Options) error {
	out = filepath.Join(out, filename)

	imageURL, err := page.ImageURL()
	if err != nil {
		return err
	}

	response, err := utils.HTTPGet(imageURL, options.HTTPRetry)
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

func cleanFilename(name string) string {
	return filenameCleanReplacer.Replace(name)
}
