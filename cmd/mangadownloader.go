package main

import (
	"fmt"
	"github.com/pierrre/mangadownloader"
)

func main() {
	var service mangadownloader.Service = &mangadownloader.MangaReaderService{}
	mangas, err := service.Mangas()
	if err != nil {
		panic(err)
	}
	for _, manga := range mangas {
		fmt.Println(manga)
		chapters, err := manga.Chapters()
		if err != nil {
			panic(err)
		}
		for _, chapter := range chapters {
			fmt.Println("	" + fmt.Sprint(chapter))
		}
	}
}
