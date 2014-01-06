package main

import (
	md "github.com/pierrre/mangadownloader"
	"github.com/pierrre/mangadownloader/service"

	"errors"
	"flag"
	"fmt"
	"net/url"
)

func main() {
	outFlag := flag.String("out", "", "Output directory")
	cbzFlag := flag.Bool("cbz", false, "CBZ")
	pageDigitCountFlag := flag.Int("pagedigitcount", 4, "Page digit count")
	httpRetryFlag := flag.Int("httpretry", 5, "Http retry")
	parallelChapterFlag := flag.Int("parallelchapter", 4, "Parallel chapter")
	parallelPageFlag := flag.Int("parallelpage", 8, "Parallel page")
	flag.Parse()
	out := *outFlag

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("Pass urls (manga/chapter/page) as argument.")
		fmt.Println("")
		fmt.Println("Options: (pass them BEFORE the arguments, Go' \"flag\" package is not really smart...)")
		flag.PrintDefaults()
		return
	}

	options := &md.Options{
		Cbz:             *cbzFlag,
		PageDigitCount:  *pageDigitCountFlag,
		ParallelChapter: *parallelChapterFlag,
		ParallelPage:    *parallelPageFlag,

		HTTPRetry: *httpRetryFlag,
	}

	for _, arg := range args {
		u, err := url.Parse(arg)
		if err != nil {
			panic(err)
		}
		o, err := md.Identify(u, options)
		if err != nil {
			panic(err)
		}
		switch object := o.(type) {
		case *service.Manga:
			err := md.DownloadManga(object, out, options)
			if err != nil {
				panic(err)
			}
		case *service.Chapter:
			err := md.DownloadChapter(object, out, options)
			if err != nil {
				panic(err)
			}
		case *service.Page:
			err := md.DownloadPage(object, out, "image", options)
			if err != nil {
				panic(err)
			}
		default:
			panic(errors.New("not supported"))
		}
	}
}
