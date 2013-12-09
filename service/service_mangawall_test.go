package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangWallTestUrlManga, _   = url.Parse("http://mangawall.com/manga/nozoki-ana")
	serviceMangWallTestUrlChapter, _ = url.Parse("http://mangawall.com/manga/nozoki-ana/6")
)

func getTestMangaWallService() *MangaWallService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangawall"]
	mangaWallService := service.(*MangaWallService)
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
