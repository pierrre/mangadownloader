package main

import (
	"errors"
	"flag"
	"github.com/pierrre/mangadownloader"
	"net/url"
)

func main() {
	outFlag := flag.String("out", "", "Output directory")
	httpRetryFlag := flag.Int("httpretry", 5, "Http retry")
	concurrencyChapterFlag := flag.Int("concurrencychapter", 4, "Concurrency chapter")
	concurrencyPageFlag := flag.Int("concurrencypage", 8, "Concurrency page")
	flag.Parse()
	out := *outFlag
	httpRetry := *httpRetryFlag
	concurrencyChapter := *concurrencyChapterFlag
	concurrencyPage := *concurrencyPageFlag

	md := mangadownloader.CreateDefaultMangeDownloader()
	md.HttpRetry = httpRetry
	md.ConcurrencyChapter = concurrencyChapter
	md.ConcurrencyPage = concurrencyPage

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
			err := md.DownloadPage(object, out)
			if err != nil {
				panic(err)
			}
		default:
			panic(errors.New("Not supported"))
		}
	}
}
