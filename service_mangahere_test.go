package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangaHereTestUrlManga, _   = url.Parse("http://www.mangahere.com/manga/aku_no_hana/")
	serviceMangaHereTestUrlChapter, _ = url.Parse("http://www.mangahere.com/manga/aku_no_hana/c012/")
)

func getTestMangaHereService() *MangaHereService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangahere"]
	mangaHereService := service.(*MangaHereService)
	return mangaHereService
}

func TestMangaHereServiceManga(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceManga(t, service, serviceMangaHereTestUrlManga)
}

func TestMangaHereServiceChapter(t *testing.T) {
	service := getTestMangaHereService()

	testCommonServiceChapter(t, service, serviceMangaHereTestUrlChapter)
}
