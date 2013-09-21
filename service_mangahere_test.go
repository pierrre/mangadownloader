package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangaHereTestUrlManga, _   = url.Parse("http://www.mangahere.com/manga/aku_no_hana/")
	serviceMangaHereTestUrlChapter, _ = url.Parse("http://www.mangahere.com/manga/aku_no_hana/c012/")
	serviceMangaHereTestUrlDummy, _   = url.Parse("http://www.google.com")
)

func getTestMangaHereService() *MangaHereService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangahere"]
	mangaHereService := service.(*MangaHereService)
	return mangaHereService
}

func TestMangaHereServiceSupports(t *testing.T) {
	t.Parallel()

	service := getTestMangaHereService()

	if !service.Supports(serviceMangaHereTestUrlManga) {
		t.Error("Not supported")
	}

	if !service.Supports(serviceMangaHereTestUrlChapter) {
		t.Error("Not supported")
	}

	if service.Supports(serviceMangaHereTestUrlDummy) {
		t.Error("Supported")
	}
}

func TestMangaHereServiceManga(t *testing.T) {
	t.Parallel()

	service := getTestMangaHereService()

	object, err := service.Identify(serviceMangaHereTestUrlManga)
	if err != nil {
		t.Fatal(err)
	}

	manga, ok := object.(*Manga)
	if !ok {
		t.Fatal("Not a manga")
	}

	name, err := manga.Name()
	if err != nil {
		t.Fatal(err)
	}
	if len(name) == 0 {
		t.Fatal("Empty name")
	}

	chapters, err := manga.Chapters()
	if err != nil {
		t.Fatal(err)
	}
	if len(chapters) == 0 {
		t.Fatal("No chapter")
	}
}

func TestMangaHereServiceChapter(t *testing.T) {
	t.Parallel()

	service := getTestMangaHereService()

	object, err := service.Identify(serviceMangaHereTestUrlChapter)
	if err != nil {
		t.Fatal(err)
	}

	chapter, ok := object.(*Chapter)
	if !ok {
		t.Fatal("Not a chapter")
	}

	name, err := chapter.Name()
	if err != nil {
		t.Fatal(err)
	}
	if len(name) == 0 {
		t.Fatal("Empty name")
	}

	pages, err := chapter.Pages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) == 0 {
		t.Fatal("No pages")
	}

	page := pages[0]
	_, err = page.ImageUrl()
	if err != nil {
		t.Fatal(err)
	}
}
