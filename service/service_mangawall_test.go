package service

import (
	"net/url"
	"testing"
)

var (
	serviceMangWallTestUrlManga, _   = url.Parse("http://mangawall.com/manga/nozoki-ana")
	serviceMangWallTestUrlChapter, _ = url.Parse("http://mangawall.com/manga/nozoki-ana/6")
)

func getTestMangaWallService() *MangaWallService {
	service := Services["mangawall"]
	mangaWallService := service.(*MangaWallService)
	mangaWallService.httpRetry = 5
	return mangaWallService
}

func TestMangaWallServiceManga(t *testing.T) {
	service := getTestMangaWallService()

	testCommonServiceManga(t, service, serviceMangWallTestUrlManga, "Nozoki Ana")
}

func TestMangaWallServiceChapter(t *testing.T) {
	service := getTestMangaWallService()

	testCommonServiceChapter(t, service, serviceMangWallTestUrlChapter, "6")
}
