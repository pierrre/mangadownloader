package main

import (
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
