package mangadownloader

type Service interface {
	Mangas() ([]*Manga, error)
}
