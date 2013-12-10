package service

import (
	"net/url"
	"testing"
)

func testCommonServiceManga(t *testing.T, service ServiceHandler, mangaURL *url.URL, expectedMangaName string) {
	//t.Parallel()

	object, err := service.Identify(mangaURL)
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
	if name != expectedMangaName {
		t.Fatal("Unexpected name")
	}

	chapters, err := manga.Chapters()
	if err != nil {
		t.Fatal(err)
	}
	if len(chapters) == 0 {
		t.Fatal("No chapter")
	}
}

func testCommonServiceChapter(t *testing.T, service ServiceHandler, chapterURL *url.URL, expectedChapterName string) {
	//t.Parallel()

	object, err := service.Identify(chapterURL)
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
	if name != expectedChapterName {
		t.Fatal("Unexpected name")
	}

	pages, err := chapter.Pages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) == 0 {
		t.Fatal("No pages")
	}

	page := pages[0]
	_, err = page.ImageURL()
	if err != nil {
		t.Fatal(err)
	}
}
