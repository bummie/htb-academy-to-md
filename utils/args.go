package utils

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	Modules         []string
	Cookies         string
	ImageDirectory  string
	OutputDirectory string
}

func GetArguements() Args {
	var moduleFlag = flag.String("m", "module", "(REQUIRED) Academy Module URL to the first page, or local file with list of module urls.")
	var cookiesFlag = flag.String("c", "cookies", "(REQUIRED) Academy Cookies for authorization.")
	var imageDirectoryFlag = flag.String("i", "imagepath", "Directory to store images locally.")
	var outputDirectoryFlag = flag.String("o", "output", "Outputdirectory to store the markdown files.")
	flag.Parse()
	arg := Args{
		Cookies:         *cookiesFlag,
		ImageDirectory:  *imageDirectoryFlag,
		OutputDirectory: *outputDirectoryFlag,
	}

	if *moduleFlag == "" || arg.Cookies == "" {
		fmt.Println("Missing required arguments for module URL and HTB Academy Cookies. Please use the -h option for help.")
		os.Exit(1)
	}

	if validateHTBUrl(*moduleFlag) {
		arg.Modules = append(arg.Modules, *moduleFlag)
	} else {
		validatePath(*moduleFlag)
		arg.Modules = readHTBUrlsFromFile(*moduleFlag)
	}

	return arg
}

func validateHTBUrl(url string) bool {
	return strings.HasPrefix(url, "https://academy.hackthebox.com/module/")
}

func validatePath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("The path '%s' does not exist\n", path)
		os.Exit(1)
	}
}

func readHTBUrlsFromFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %s", err.Error())
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if validateHTBUrl(line) {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s", err.Error())
		os.Exit(2)
	}

	return lines
}
