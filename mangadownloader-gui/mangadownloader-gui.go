package main

import (
	"github.com/matrixik/mangadownloader"
	"github.com/matrixik/mangadownloader/service"

	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

const (
	HTTP_ADDR_DEFAULT = "localhost:0"
)

var (
	md *mangadownloader.MangaDownloader
)

func init() {
	md = mangadownloader.CreateDefaultMangeDownloader()
}

func main() {
	httpAddrFlag := flag.String("http", HTTP_ADDR_DEFAULT, "Http")
	openBrowserFlag := flag.Bool("browser", true, "Browser")
	flag.Parse()
	httpAddr := *httpAddrFlag
	openBrowser := *openBrowserFlag

	listener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		panic(err)
	}

	if openBrowser {
		u := &url.URL{
			Scheme: "http",
			Host:   listener.Addr().String(),
		}

		err = openUrl(u)
		if err != nil {
			log.Println(err)
			log.Printf("Open '%s' manually", u)
		}
	}

	http.HandleFunc("/", httpHandleIndex)
	http.HandleFunc("/add", httpHandleAdd)

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}

func openUrl(u *url.URL) error {
	us := u.String()
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", us)
	case "windows":
		command = exec.Command("cmd", "/c", "start", us)
	default:
		command = exec.Command("xdg-open", us)
	}
	return command.Run()
}

func httpHandleIndex(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	writer.Write([]byte(`<html>
	<head>
		<title>Manga Downloader</title>
	</head>
	<body>
		<form action="/add" method="post">
			<input type="text" name="url"><br />
			<input type="submit" value="Download">
		</form>
	</body>
</html>`))
}

func httpHandleAdd(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	err := request.ParseForm()
	if err != nil {
		httpError(writer, err)
		return
	}

	u, err := url.Parse(request.FormValue("url"))
	if err != nil {
		httpError(writer, err)
		return
	}

	go func() {
		o, err := md.Identify(u)
		if err != nil {
			return
		}
		options := &mangadownloader.Options{
			Cbz:             true,
			PageDigitCount:  4,
			ParallelChapter: 4,
			ParallelPage:    8,
		}
		out := ""
		switch object := o.(type) {
		case *service.Manga:
			md.DownloadManga(object, out, options)
		case *service.Chapter:
			md.DownloadChapter(object, out, options)
		case *service.Page:
			md.DownloadPage(object, out, "image", options)
		}
	}()
}

func httpError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(err.Error()))
}
