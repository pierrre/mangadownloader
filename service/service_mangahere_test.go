package service

import (
	"net/url"
	"testing"
)

var (
	serviceMangaHereTestURLManga, _   = url.Parse("http://www.mangahere.com/manga/aku_no_hana/")
	serviceMangaHereTestURLChapter, _ = url.Parse("http://www.mangahere.com/manga/aku_no_hana/c012/")
)

func getTestMangaHereService() *MangaHereService {
	service := Services["mangahere"]
	mangaHereService := service.(*MangaHereService)
	mangaHereService.httpRetry = 5
	return mangaHereService
}

func TestMangaHereServiceManga(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceManga(t, service, serviceMangaHereTestURLManga, "Aku No Hana")
}

func TestMangaHereServiceChapter(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceChapter(t, service, serviceMangaHereTestURLChapter, "012")
}
