package main

import (
	"errors"
	"flag"
	"github.com/pierrre/mangadownloader"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	out, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	md := mangadownloader.CreateDefaultMangeDownloader()

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
		default:
			panic(errors.New("Not supported"))
		}
	}
}
