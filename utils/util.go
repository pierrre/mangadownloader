package utils

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"errors"
	"net/http"
	"net/url"
	"os"
)

func HTTPGetHTML(u *url.URL, HTTPRetry int) (*html.Node, error) {
	response, err := HTTPGet(u, HTTPRetry)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	node, err := html.Parse(response.Body)
	return node, err
}

func HTTPGet(u *url.URL, HTTPRetry int) (response *http.Response, err error) {
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.4 Safari/537.36")

	httpRetry := HTTPRetry
	if httpRetry < 1 {
		httpRetry = 1
	}

	errs := make(MultiError, 0)
	for i := 0; i < httpRetry; i++ {
		response, err := http.DefaultClient.Do(request)
		if err == nil {
			return response, nil
		}
		errs = append(errs, err)
	}
	return nil, errs
}

func HTMLGetNodeAttribute(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func HTMLGetNodeText(node *html.Node) (string, error) {
	switch node.Type {
	case html.TextNode:
		return node.Data, nil
	case html.DocumentNode, html.ElementNode:
		buffer := new(bytes.Buffer)
		childNode := node.FirstChild
		for childNode != nil {
			text, err := HTMLGetNodeText(childNode)
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

func URLCopy(u *url.URL) *url.URL {
	urlCopyVal := *u
	urlCopy := &urlCopyVal
	return urlCopy
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func StringSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
