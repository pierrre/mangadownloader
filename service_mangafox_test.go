package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangaFoxTestUrlManga, _   = url.Parse("http://mangafox.me/manga/berserk/")
	serviceMangaFoxTestUrlChapter, _ = url.Parse("http://mangafox.me/manga/berserk/c134/1.html")
)

func getTestMangaFoxService() *MangaFoxService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangafox"]
	mangaFoxService := service.(*MangaFoxService)
	return mangaFoxService
}

func TestMangaFoxServiceManga(t *testing.T) {
	service := getTestMangaFoxService()

	testCommonServiceManga(t, service, serviceMangaFoxTestUrlManga)
}

func TestMangaFoxServiceChapter(t *testing.T) {
	service := getTestMangaFoxService()

	testCommonServiceChapter(t, service, serviceMangaFoxTestUrlChapter)
}
