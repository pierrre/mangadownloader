package mangadownloader

import (
	"net/url"
	"testing"
)

var (
	serviceMangReaderTestUrlManga, _   = url.Parse("http://www.mangareader.net/97/gantz.html")
	serviceMangReaderTestUrlChapter, _ = url.Parse("http://www.mangareader.net/97-1220-12/gantz/chapter-58.html")
	serviceMangReaderTestUrlDummy, _   = url.Parse("http://www.google.com")
)

func getTestMangaReaderService() *MangaReaderService {
	md := CreateDefaultMangeDownloader()
	for _, service := range md.Services {
		if mangaReaderService, ok := service.(*MangaReaderService); ok {
			return mangaReaderService
		}
	}
	return nil
}

func TestSupports(t *testing.T) {
	t.Parallel()

	service := getTestMangaReaderService()

	if !service.Supports(serviceMangReaderTestUrlManga) {
		t.Error("Not supported")
	}

	if !service.Supports(serviceMangReaderTestUrlChapter) {
		t.Error("Not supported")
	}

	if service.Supports(serviceMangReaderTestUrlDummy) {
		t.Error("Supported")
	}
}

func TestMangas(t *testing.T) {
	t.Parallel()

	service := getTestMangaReaderService()

	mangas, err := service.Mangas()
	if err != nil {
		t.Fatal(err)
	}

	if len(mangas) == 0 {
		t.Fatal("No manga")
	}
}

func TestManga(t *testing.T) {
	t.Parallel()

	service := getTestMangaReaderService()

	object, err := service.Identify(serviceMangReaderTestUrlManga)
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

func TestChapter(t *testing.T) {
	t.Parallel()

	service := getTestMangaReaderService()

	object, err := service.Identify(serviceMangReaderTestUrlChapter)
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
