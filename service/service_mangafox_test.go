package service

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
	if testing.Short() {
		t.Skip("it fails randomly on Travis CI")
	}

	service := getTestMangaFoxService()

	testCommonServiceManga(t, service, serviceMangaFoxTestUrlManga, "Berserk")
}

func TestMangaFoxServiceChapter(t *testing.T) {
	if testing.Short() {
		t.Skip("it fails randomly on Travis CI")
	}

	service := getTestMangaFoxService()

	testCommonServiceChapter(t, service, serviceMangaFoxTestUrlChapter, "134")
}
