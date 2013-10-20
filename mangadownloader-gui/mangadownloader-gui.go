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
	HTTP_ADDR_DEFAULT = ":0"
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
	commandName := openUrlCommandName()
	command := exec.Command(commandName, u.String())
	return command.Run()
}

func openUrlCommandName() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	case "windows":
		return "start"
	default:
		return "xdg-open"
	}
}
