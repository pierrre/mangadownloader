package service

import (
	"net/url"
	"testing"
)

var (
	serviceTenMangaTestUrlManga, _   = url.Parse("http://www.tenmanga.com/book/Green+Blood.html")
	serviceTenMangaTestUrlChapter, _ = url.Parse("http://www.tenmanga.com/chapter/GreenBlood9/299502/")
)

func getTestTenMangaService() *TenMangaService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["tenmanga"]
	tenMangaService := service.(*TenMangaService)
	return tenMangaService
}

func TestTenMangaServiceManga(t *testing.T) {
	service := getTestTenMangaService()

	testCommonServiceManga(t, service, serviceTenMangaTestUrlManga, "Green Blood")
}

func TestTenMangaServiceChapter(t *testing.T) {
	service := getTestTenMangaService()

	testCommonServiceChapter(t, service, serviceTenMangaTestUrlChapter, "9")
}
