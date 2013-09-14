package main

import (
	"errors"
	"flag"
	"github.com/pierrre/mangadownloader"
	"net/url"
)

func main() {
	outFlag := flag.String("out", "", "Output directory")
	pageDigitCountFlag := flag.Int("pagedigitcount", 4, "Page digit count")
	httpRetryFlag := flag.Int("httpretry", 5, "Http retry")
	parallelChapterFlag := flag.Int("parallelchapter", 4, "Parallel chapter")
	parallelPageFlag := flag.Int("parallelpage", 8, "Parallel page")
	flag.Parse()
	out := *outFlag

	md := mangadownloader.CreateDefaultMangeDownloader()
	md.PageDigitCount = *pageDigitCountFlag
	md.HttpRetry = *httpRetryFlag
	md.ParallelChapter = *parallelChapterFlag
	md.ParallelPage = *parallelPageFlag

	for _, arg := range flag.Args() {
		u, err := url.Parse(arg)
		if err != nil {
			panic(err)
		}
		o, err := md.Identify(u)
		if err != nil {
			panic(err)
		}
		switch object := o.(type) {
		case *mangadownloader.Manga:
			err := md.DownloadManga(object, out)
			if err != nil {
				panic(err)
			}
		case *mangadownloader.Chapter:
			err := md.DownloadChapter(object, out)
			if err != nil {
				panic(err)
			}
		case *mangadownloader.Page:
			err := md.DownloadPage(object, out, "image")
			if err != nil {
				panic(err)
			}
		default:
			panic(errors.New("Not supported"))
		}
	}
}
