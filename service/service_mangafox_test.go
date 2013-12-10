package service

import (
	"net/url"
	"testing"
)

var (
	serviceMangaFoxTestURLManga, _   = url.Parse("http://mangafox.me/manga/berserk/")
	serviceMangaFoxTestURLChapter, _ = url.Parse("http://mangafox.me/manga/berserk/c134/1.html")
)

func getTestMangaFoxService() *MangaFoxService {
	service := Services["mangafox"]
	mangaFoxService := service.(*MangaFoxService)
	mangaFoxService.httpRetry = 5
	return mangaFoxService
}

func TestMangaFoxServiceManga(t *testing.T) {
	if testing.Short() {
		t.Skip("it fails randomly on Travis CI")
	}

	service := getTestMangaFoxService()

	testCommonServiceManga(t, service, serviceMangaFoxTestURLManga, "Berserk")
}

func TestMangaFoxServiceChapter(t *testing.T) {
	if testing.Short() {
		t.Skip("it fails randomly on Travis CI")
	}

	service := getTestMangaFoxService()

	testCommonServiceChapter(t, service, serviceMangaFoxTestURLChapter, "134")
}
