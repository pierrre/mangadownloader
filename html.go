package mangadownloader

import (
	"code.google.com/p/go.net/html"
)

func getHtmlNodeAttribute(node *html.Node, key string) (value string) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
