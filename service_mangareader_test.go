package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangReaderTestUrlManga, _   = url.Parse("http://www.mangareader.net/97/gantz.html")
	serviceMangReaderTestUrlChapter, _ = url.Parse("http://www.mangareader.net/97-1220-12/gantz/chapter-58.html")
)

func getTestMangaReaderService() *MangaReaderService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangareader"]
	mangaReaderService := service.(*MangaReaderService)
	return mangaReaderService
}

func TestMangaReaderServiceManga(t *testing.T) {
	service := getTestMangaReaderService()

	testCommonServiceManga(t, service, serviceMangReaderTestUrlManga)
}

func TestMangaReaderServiceChapter(t *testing.T) {
	service := getTestMangaReaderService()

	testCommonServiceChapter(t, service, serviceMangReaderTestUrlChapter)
}
