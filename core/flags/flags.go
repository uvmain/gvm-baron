package flags

import (
	"flag"
	"log"
	"os"
	"slices"
	"strings"
)

var ListType string
var ListEnabled bool

var DownloadEnabled bool
var DownloadVersion string

var DebugEnabled bool

func InitFlags() {
	args := os.Args[1:]
	if slices.Contains(args, "--list") || slices.Contains(args, "-list") {
		ListEnabled = true
	}
	for i, arg := range args {
		if arg == "--list" || arg == "-list" {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				args[i] = arg + "=stable"
			}
		}
	}
	os.Args = append([]string{os.Args[0]}, args...)

	if slices.Contains(args, "--download") || slices.Contains(args, "-download") {
		DownloadEnabled = true
	}

	listType := flag.String("list", "stable", "latest | stable | lts | all")

	flag.BoolVar(&DebugEnabled, "debug", false, "Enable debug mode")

	downloadVersion := flag.String("download", "", "Specify the version to download")

	flag.Parse()

	ListType = *listType
	DownloadVersion = *downloadVersion

	if ListEnabled && !slices.Contains([]string{"latest", "stable", "lts", "all"}, ListType) {
		log.Printf("Invalid list type: %s. Valid options are: latest, stable, lts, all. Defaulting to 'stable'", ListType)
		ListType = "stable"
	}
}
