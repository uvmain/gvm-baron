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

var InstallEnabled bool
var InstallVersion string

var AliasEnabled bool
var AliasSource string
var AliasTarget string

var DebugEnabled bool
var NoCache bool

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

	if slices.Contains(args, "--install") || slices.Contains(args, "-install") {
		InstallEnabled = true
	}

	listType := flag.String("list", "stable", "latest | stable | lts | all")

	flag.BoolVar(&DebugEnabled, "debug", false, "Enable debug mode")
	flag.BoolVar(&NoCache, "no-cache", false, "Disable cache")
	flag.BoolVar(&AliasEnabled, "alias", false, "Create alias for a Go version")
	flag.StringVar(&AliasSource, "alias-source", "", "Source version for alias (e.g., 1.20)")
	flag.StringVar(&AliasTarget, "alias-target", "", "Target alias name (e.g., latest)")
	installVersion := flag.String("install", "", "latest | lts | x.y.z | x.y")

	flag.Parse()

	ListType = *listType
	InstallVersion = *installVersion

	if ListEnabled && !slices.Contains([]string{"latest", "stable", "lts", "all"}, ListType) {
		log.Printf("Invalid list type: %s. Valid options are: latest, stable, lts, all. Defaulting to 'stable'", ListType)
		ListType = "stable"
	}
}
