package service

import (
	"net/url"
	"testing"
)

var (
	serviceMangWallTestURLManga, _   = url.Parse("http://mangawall.com/manga/nozoki-ana")
	serviceMangWallTestURLChapter, _ = url.Parse("http://mangawall.com/manga/nozoki-ana/6")
)

func getTestMangaWallService() *MangaWallService {
	service := Services["mangawall"]
	mangaWallService := service.(*MangaWallService)
	mangaWallService.httpRetry = 5
	return mangaWallService
}

func TestMangaWallServiceManga(t *testing.T) {
	service := getTestMangaWallService()

	testCommonServiceManga(t, service, serviceMangWallTestURLManga, "Nozoki Ana")
}

func TestMangaWallServiceChapter(t *testing.T) {
	service := getTestMangaWallService()

	testCommonServiceChapter(t, service, serviceMangWallTestURLChapter, "6")
}
