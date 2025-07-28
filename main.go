package main

import (
	"fmt"
	"htb-academy-md/parser"
	"htb-academy-md/utils"
	"htb-academy-md/webrequest"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func main() {
	options := utils.GetArguements()

	fmt.Println("Authenticating with HackTheBox...")
	session, err := webrequest.AuthenticateWithCookies(options.Cookies)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed authencating to HackTheBox: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("Downloading requested module...")
	title, content, err := parser.GetModule(options.ModuleUrl, session)

	if options.LocalImages {
		fmt.Println("Downloading module images...")
		content, err = parser.GetImagesLocally(content)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed downloading images locally %s", err.Error())
			os.Exit(4)
		}
	}

	markdownContent, err := htmlToMarkdown(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed converting html document to markdown %s", err.Error())
		os.Exit(2)
	}

	err = os.WriteFile(title+".md", []byte(markdownContent), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed writing markdown content to file %s", err.Error())
		os.Exit(3)
	}

	fmt.Println("Finished downloading module!")
}

func htmlToMarkdown(html []string) (string, error) {
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	var markdown string
	for _, content := range html {
		m, err := converter.ConvertString(content)
		if err != nil {
			return "", err
		}
		markdown += m + "\n\n\n"
	}

	// Strip some content for proper code blocks.
	markdown = strings.ReplaceAll(markdown, "shell-session", "shell")
	markdown = strings.ReplaceAll(markdown, "powershell-session", "powershell")
	markdown = strings.ReplaceAll(markdown, "[!bash!]$ ", "")

	return markdown, nil
}
