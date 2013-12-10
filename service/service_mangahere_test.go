package service

import (
	"net/url"
	"testing"
)

var (
	serviceMangaHereTestUrlManga, _   = url.Parse("http://www.mangahere.com/manga/aku_no_hana/")
	serviceMangaHereTestUrlChapter, _ = url.Parse("http://www.mangahere.com/manga/aku_no_hana/c012/")
)

func getTestMangaHereService() *MangaHereService {
	service := Services["mangahere"]
	mangaHereService := service.(*MangaHereService)
	mangaHereService.httpRetry = 5
	return mangaHereService
}

func TestMangaHereServiceManga(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceManga(t, service, serviceMangaHereTestUrlManga, "Aku No Hana")
}

func TestMangaHereServiceChapter(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceChapter(t, service, serviceMangaHereTestUrlChapter, "012")
}
