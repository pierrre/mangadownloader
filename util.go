package mangadownloader

import (
	"code.google.com/p/go.net/html"
	"net/url"
	"os"
)

func htmlGetNodeAttribute(node *html.Node, key string) (value string) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func urlCopy(u *url.URL) *url.URL {
	urlCopyVal := *u
	urlCopy := &urlCopyVal
	return urlCopy
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
