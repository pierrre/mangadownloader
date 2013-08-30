package mangadownloader

type Service interface {
	Mangas() ([]*Manga, error)
	Chapters(*Manga) ([]*Chapter, error)
}
