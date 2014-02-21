package mangadownloader

import (
	"bytes"
	"errors"
	"net/url"
	"os"

	"code.google.com/p/go.net/html"
)

func htmlGetNodeAttribute(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func htmlGetNodeText(node *html.Node) (string, error) {
	switch node.Type {
	case html.TextNode:
		return node.Data, nil
	case html.DocumentNode, html.ElementNode:
		buffer := new(bytes.Buffer)
		childNode := node.FirstChild
		for childNode != nil {
			text, err := htmlGetNodeText(childNode)
			if err != nil {
				return "", err
			}
			_, err = buffer.WriteString(text)
			if err != nil {
				return "", err
			}
			childNode = childNode.NextSibling
		}
		return buffer.String(), nil
	case html.CommentNode:
		return "", nil
	case html.DoctypeNode:
		return "", nil
	case html.ErrorNode:
		return "", nil
	default:
		return "", errors.New("invalid html node type")
	}
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

func stringSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func chapterSliceReverse(chapters []*Chapter) []*Chapter {
	count := len(chapters)
	reversed := make([]*Chapter, 0, count)
	for i := count - 1; i >= 0; i-- {
		reversed = append(reversed, chapters[i])
	}
	return reversed
}
