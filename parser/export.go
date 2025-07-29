package parser

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func GetModule(moduleUrl string, client *http.Client) (string, []string, error) {
	resp, err := client.Get(moduleUrl)
	if err != nil {
		return "", []string{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", []string{}, err
	}

	content := string(body)
	moduleTitle, err := getModuleTitle(content)
	if err != nil {
		return "", []string{}, err
	}

	pageUrls, err := getModulePages(content, moduleUrl)
	if err != nil {
		return "", []string{}, err
	}

	var pagesContent []string
	for _, pageUrl := range pageUrls {
		content, err := extractPageContent(pageUrl, client)
		if err != nil {
			return "", []string{}, err
		}

		pagesContent = append(pagesContent, content)
	}

	return moduleTitle, pagesContent, nil
}

func getModuleTitle(htmlText string) (string, error) {
	var title string
	var isTitle bool
	badChars := []string{"/", "\\", "?", "%", "*", ":", "|", "\"", "<", ">"}
	tkn := html.NewTokenizer(strings.NewReader(htmlText))

	for {
		tt := tkn.Next()
		switch {

		case tt == html.ErrorToken:
			os.Exit(1)

		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "title" {
				isTitle = true
			}
		case tt == html.TextToken:
			t := tkn.Token()

			if isTitle {
				title = t.Data
				for _, badChar := range badChars {
					title = strings.ReplaceAll(title, badChar, "-")
				}
				return title, nil
			}
		case tt == html.EndTagToken:
			t := tkn.Token()
			if t.Data == "html" {
				return "", errors.New("Could not find title on module page, exiting.")
			}
		}
	}
}

func getModulePages(htmlText string, moduleUrl string) ([]string, error) {
	var modulePages []string

	doc, err := html.Parse(strings.NewReader(htmlText))
	if err != nil {
		return nil, err
	}

	var traverse func(n *html.Node) *html.Node
	traverse = func(n *html.Node) *html.Node {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "a" {
				// Pass on <a> tags that do not have any attributes.
				if len(c.Attr) == 0 {
				} else if strings.Contains(c.Attr[0].Val, moduleUrl[:44]) {
					modulePages = append(modulePages, c.Attr[0].Val)
				}
			}
			res := traverse(c)
			if res != nil {
				return res
			}
		}
		return nil
	}
	traverse(doc)

	return modulePages[1:], nil
}

func getModulePageContent(pageUrl string, client *http.Client) (string, error) {
	resp, err := client.Get(pageUrl)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	content := string(body)
	return content, nil
}

func extractPageContent(pageUrl string, client *http.Client) (string, error) {
	var result string
	pageContent, err := getModulePageContent(pageUrl, client)

	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	trainingContent := findDivByClass(doc, "training-module")
	if trainingContent != nil {
		var tc strings.Builder
		if err := html.Render(&tc, trainingContent); err != nil {
			return "", err
		}
		result = tc.String()
	} else {
		return "", errors.New(fmt.Sprintf("Parsing training content failed, HTML dump: %s", pageContent))
	}

	currentDoc, err := html.Parse(strings.NewReader(result))
	if err != nil {
		return "", nil
	}

	contentToRemove := findDivByClassOrId(currentDoc, "vpn-switch-card", "screen")
	if contentToRemove != nil {
		var htmlBuilder strings.Builder
		if err := html.Render(&htmlBuilder, contentToRemove); err != nil {
			return "", err
		}

		indexToRemove := strings.Index(result, htmlBuilder.String())
		result = result[:indexToRemove]
	}

	curDoc, err := html.Parse(strings.NewReader(result))
	if err != nil {
		return "", err
	}

	fixImgs(curDoc)
	var r strings.Builder
	if err := html.Render(&r, curDoc); err != nil {
		return "", err
	}
	result = r.String()

	return result, nil
}

func findDivByClass(n *html.Node, className string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, className) {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if div := findDivByClass(c, className); div != nil {
			return div
		}
	}

	return nil
}

func findDivByClassOrId(n *html.Node, className string, id string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if (a.Key == "class" && strings.Contains(a.Val, className)) ||
				(a.Key == "id" && strings.Contains(a.Val, id)) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if div := findDivByClassOrId(c, className, id); div != nil {
			return div
		}
	}

	return nil
}

func fixImgs(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "img" {
		for i, attr := range node.Attr {
			if attr.Key == "src" && !startsWith(attr.Val, "https://") {
				node.Attr[i].Val = "https://academy.hackthebox.com" + attr.Val
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		fixImgs(child)
	}
}

func startsWith(str, prefix string) bool {
	return len(str) >= len(prefix) && str[0:len(prefix)] == prefix
}

func GetImagesLocally(htmlPages []string, writePath string) ([]string, error) {
	var result []string
	for _, page := range htmlPages {
		doc, err := html.Parse(strings.NewReader(page))
		if err != nil {
			return []string{}, err
		}

		replaceImgs(doc, writePath)
		var htmlBulder strings.Builder
		if err := html.Render(&htmlBulder, doc); err != nil {
			return []string{}, err
		}

		result = append(result, htmlBulder.String())
	}

	return result, nil
}

func replaceImgs(node *html.Node, writePath string) {
	if node.Type == html.ElementNode && node.Data == "img" {
		for i, attr := range node.Attr {
			if attr.Key == "src" {
				fileName, err := downloadImage(node.Attr[i].Val, writePath)
				if err != nil {
					// Replace url, but send out message that the image was not correctly downloaded
					fmt.Fprintln(os.Stderr, "could not download image "+node.Attr[i].Val)
				}
				node.Attr[i].Val = fileName
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		replaceImgs(child, writePath)
	}
}

func downloadImage(fileUrl string, writePath string) (string, error) {
	fileName := randomFileName()
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if isPNG(content) {
		fileName = fileName + ".png"
		fileName = filepath.Join(writePath, fileName)
		err := os.WriteFile(fileName, content, 0666)
		if err != nil {
			return "", err
		}
	} else if isJPEG(content) {
		fileName = fileName + ".jpg"
		fileName = filepath.Join(writePath, fileName)
		err := os.WriteFile(fileName, content, 0666)
		if err != nil {
			return "", err
		}
	} else if isGIF(content) {
		fileName = fileName + ".gif"
		fileName = filepath.Join(writePath, fileName)
		err := os.WriteFile(fileName, content, 0666)
		if err != nil {
			return "", err
		}
	} else {
		fileName = filepath.Join(writePath, fileName)
		err := os.WriteFile(fileName, content, 0666)
		if err != nil {
			return "", err
		}
	}

	return fileName, nil
}

func randomFileName() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 12)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func isPNG(b []byte) bool {
	return len(b) >= 8 && string(b[0:8]) == "\x89PNG\r\n\x1a\n"
}

func isJPEG(b []byte) bool {
	return len(b) >= 2 && string(b[0:2]) == "\xff\xd8"
}

func isGIF(data []byte) bool {
	if len(data) < 6 {
		return false
	}
	if string(data[:3]) != "GIF" {
		return false
	}
	if string(data[3:6]) != "87a" && string(data[3:6]) != "89a" {
		return false
	}
	return true
}
