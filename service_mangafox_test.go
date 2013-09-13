package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangFoxTestUrlManga, _   = url.Parse("http://mangafox.me/manga/berserk/")
	serviceMangFoxTestUrlChapter, _ = url.Parse("http://mangafox.me/manga/berserk/c134/1.html")
	serviceMangFoxTestUrlDummy, _   = url.Parse("http://www.google.com")
)

func getTestMangaFoxService() *MangaFoxService {
	md := CreateDefaultMangeDownloader()
	for _, service := range md.Services {
		if mangaFoxService, ok := service.(*MangaFoxService); ok {
			return mangaFoxService
		}
	}
	return nil
}

func TestMangaFoxServiceSupports(t *testing.T) {
	t.Parallel()

	service := getTestMangaFoxService()

	if !service.Supports(serviceMangFoxTestUrlManga) {
		t.Error("Not supported")
	}

	if !service.Supports(serviceMangFoxTestUrlChapter) {
		t.Error("Not supported")
	}

	if service.Supports(serviceMangFoxTestUrlDummy) {
		t.Error("Supported")
	}
}

func TestMangaFoxServiceManga(t *testing.T) {
	t.Parallel()

	service := getTestMangaFoxService()

	object, err := service.Identify(serviceMangFoxTestUrlManga)
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

	object, err := service.Identify(serviceMangFoxTestUrlChapter)
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
