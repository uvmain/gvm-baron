package flags

import (
	"flag"
	"os"
	"slices"
	"strings"
)

var ListType string
var ListEnabled bool

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

	listType := flag.String("list", "stable", "latest | stable | lts | all")
	flag.Parse()
	ListType = *listType
}
