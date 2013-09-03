package mangadownloader

type Service interface {
	Mangas() ([]*Manga, error)
	MangaName(*Manga) (string, error)
	Chapters(*Manga) ([]*Chapter, error)
	Pages(*Chapter) ([]*Page, error)
	Image(*Page) (*Image, error)
}
