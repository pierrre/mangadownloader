package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangaFoxTestUrlManga, _   = url.Parse("http://mangafox.me/manga/bokurano/")
	serviceMangaFoxTestUrlChapter, _ = url.Parse("http://mangafox.me/manga/bokurano/c010/1.html")
)

func getTestMangaFoxService() *MangaFoxService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangafox"]
	mangaFoxService := service.(*MangaFoxService)
	return mangaFoxService
}

func TestMangaFoxServiceManga(t *testing.T) {
	service := getTestMangaFoxService()

	testCommonServiceManga(t, service, serviceMangaFoxTestUrlManga, "Bokurano")
}

func TestMangaFoxServiceChapter(t *testing.T) {
	service := getTestMangaFoxService()

	testCommonServiceChapter(t, service, serviceMangaFoxTestUrlChapter, "010")
}
