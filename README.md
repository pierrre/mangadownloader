# MangaDownloader
A manga downloader written in Go (Golang).

## Features
- Command line
- Download manga from websites (more to come)
    - http://www.mangareader.net
    - http://mangafox.me
- Parallel download

## Binaries
https://www.dropbox.com/sh/zc94728pgccei17/mskEqE4XuM

## Usage
`./mangadownloader [urls...]`

Examples:

```
./mangadownloader http://www.mangareader.net/97/gantz.html
./mangadownloader http://mangafox.me/manga/berserk/c134/1.html
```

Options:

`./mangadownloader -h`

```
Usage of ./mangadownloader:
  -httpretry=5: Http retry
  -out="": Output directory
  -pagedigitcount=4: Page digit count
  -parallelchapter=4: Parallel chapter
  -parallelcypage=8: Parallel page
```

## Build
`go build cmd/mangadownloader.go`

## TODO
- Readme
- More services
    - animea
    - unixmanga
    - goodmanga
    - mangacraze
    - mangago
    - anymanga
    - mangainn
    - mangaeden
    - mangable
    - deliciousmanga
    - mangahere
    - tenmanga
    - mangawall
    - mangastream
    - fakku
    - doujin-moe
- Tests
- Documentation
- Sync
- Input file
- Concurrency for input
- Progress
- User agent
- Improve error handling
- Improve http error codes