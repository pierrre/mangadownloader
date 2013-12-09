package service

import (
	"net/url"
	"testing"
)

func testCommonServiceManga(t *testing.T, service ServiceHandler, mangaUrl *url.URL, expectedMangaName string) {
	//t.Parallel()

	object, err := service.Identify(mangaUrl)
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

func testCommonServiceChapter(t *testing.T, service Service, chapterUrl *url.URL, expectedChapterName string) {
	//t.Parallel()

	object, err := service.Identify(chapterUrl)
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
	_, err = page.ImageUrl()
	if err != nil {
		t.Fatal(err)
	}
}
