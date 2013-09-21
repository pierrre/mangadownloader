package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangaFoxTestUrlManga, _   = url.Parse("http://mangafox.me/manga/berserk/")
	serviceMangaFoxTestUrlChapter, _ = url.Parse("http://mangafox.me/manga/berserk/c134/1.html")
	serviceMangaFoxTestUrlDummy, _   = url.Parse("http://www.google.com")
)

func getTestMangaFoxService() *MangaFoxService {
	md := CreateDefaultMangeDownloader()
	service := md.Services["mangafox"]
	mangaFoxService := service.(*MangaFoxService)
	return mangaFoxService
}

func TestMangaFoxServiceSupports(t *testing.T) {
	t.Parallel()

	service := getTestMangaFoxService()

	if !service.Supports(serviceMangaFoxTestUrlManga) {
		t.Error("Not supported")
	}

	if !service.Supports(serviceMangaFoxTestUrlChapter) {
		t.Error("Not supported")
	}

	if service.Supports(serviceMangaFoxTestUrlDummy) {
		t.Error("Supported")
	}
}

func TestMangaFoxServiceManga(t *testing.T) {
	t.Parallel()

	service := getTestMangaFoxService()

	object, err := service.Identify(serviceMangaFoxTestUrlManga)
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

func TestMangaFoxServiceChapter(t *testing.T) {
	t.Parallel()

	service := getTestMangaFoxService()

	object, err := service.Identify(serviceMangaFoxTestUrlChapter)
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
