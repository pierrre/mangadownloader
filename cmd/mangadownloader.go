package main

import (
	"fmt"
	"github.com/pierrre/mangadownloader"
	"net/url"
)

func main() {
	url, err := url.Parse("http://www.mangareader.net/")
	if err != nil {
		panic(err)
	}
	service := &mangadownloader.Service{
		Url: url,
	}
	mangas, err := service.Mangas()
	if err != nil {
		panic(err)
	}
	for _, manga := range mangas {
		fmt.Println(manga)
	}
}
