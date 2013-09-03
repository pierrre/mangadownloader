package main

import (
	"fmt"
	"github.com/pierrre/mangadownloader"
)

func main() {
	md := mangadownloader.CreateDefaultMangeDownloader()
	for _, service := range md.Services {
		mangas, err := service.Mangas()
		if err != nil {
			panic(err)
		}
		for _, manga := range mangas {
			fmt.Println(manga)
			mangaName, err := manga.Name()
			if err != nil {
				panic(err)
			}
			fmt.Println(mangaName)
			chapters, err := manga.Chapters()
			if err != nil {
				panic(err)
			}
			for _, chapter := range chapters {
				fmt.Println("	" + fmt.Sprint(chapter))
				pages, err := chapter.Pages()
				if err != nil {
					panic(err)
				}
				for _, page := range pages {
					fmt.Println("		" + fmt.Sprint(page))
					pageIndex, err := page.Index()
					if err != nil {
						panic(err)
					}
					fmt.Println("		" + fmt.Sprint(pageIndex))
					image, err := page.Image()
					if err != nil {
						panic(err)
					}
					fmt.Println("		" + fmt.Sprint(image))
				}
			}
		}
	}
}
