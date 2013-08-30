package main

import (
	"fmt"
	"github.com/pierrre/mangadownloader"
)

func main() {
	service := &mangadownloader.MangaReaderService{}
	mangas, err := service.Mangas()
	if err != nil {
		panic(err)
	}
	for _, manga := range mangas {
		fmt.Println(manga)
	}
}
