package mangadownloader

type MangaDownloader struct {
	Services []Service
}

func CreateDefaultMangeDownloader() *MangaDownloader {
	mangaDownloader := &MangaDownloader{}

	mangaDownloader.Services = append(mangaDownloader.Services, &MangaReaderService{})

	return mangaDownloader
}
