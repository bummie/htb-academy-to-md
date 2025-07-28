package utils

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	ModuleUrl   string
	Cookies     string
	LocalImages bool
}

func GetArguements() Args {
	var mFlag = flag.String("m", "", "(REQUIRED) Academy Module URL to the first page.")
	var cFlag = flag.String("c", "", "(REQUIRED) Academy Cookies for authorization.")
	var imgFlag = flag.Bool("local_images", false, "Save images locally rather than referencing the URL location.")
	flag.Parse()
	arg := Args{
		ModuleUrl:   *mFlag,
		Cookies:     *cFlag,
		LocalImages: *imgFlag,
	}

	if arg.ModuleUrl == "" || arg.Cookies == "" {
		fmt.Println("Missing required arguments for module URL and HTB Academy Cookies. Please use the -h option for help.")
		os.Exit(1)
	}

	return arg
}
