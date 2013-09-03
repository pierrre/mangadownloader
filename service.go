package mangadownloader

type Service interface {
	Mangas() ([]*Manga, error)
	MangaName(*Manga) (string, error)
	Chapters(*Manga) ([]*Chapter, error)
	Pages(*Chapter) ([]*Page, error)
	PageIndex(*Page) (uint, error)
	Image(*Page) (*Image, error)
}
