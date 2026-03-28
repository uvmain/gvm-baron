package main

import (
	"fmt"
	"gvm/core/aliases"
	"gvm/core/config"
	"gvm/core/installed"
	"gvm/core/logger"
	"gvm/core/sources"
	"os"
	"slices"
	"strings"
)

func main() {
	config.InitConfig()

	arguments := os.Args[1:]

	args := make([]string, 0)
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		}
	}

	flags := make([]string, 0)
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		}
	}

	if len(args) == 0 || slices.Contains(flags, "--help") || slices.Contains(flags, "-h") {
		if len(arguments) == 0 {
			fmt.Println("No action specified. Valid actions are: list, install, alias.")
			fmt.Println()
		}
		printHelp()

		currentVersion, _ := aliases.GetCurrentDefaultVersion()
		if currentVersion != "" {
			fmt.Println()
			fmt.Printf("Current Go version: %s\n\n", currentVersion)
		}
		return
	}

	if slices.Contains(flags, "--debug") || slices.Contains(flags, "-d") {
		logger.DebugPrintln("Debug mode enabled")
		logger.DebugEnabled = true
	}

	if slices.Contains(flags, "--no-cache") || slices.Contains(flags, "-n") {
		logger.DebugPrintln("Cache disabled")
		config.NoCache = true
	} else {
		logger.DebugPrintln("Cache enabled")
	}

	action := args[0]

	switch action {
	case "current":
		logger.DebugPrintln("Getting current Go version...")
		currentVersion, err := aliases.GetCurrentDefaultVersion()
		if err != nil {
			fmt.Printf("Error getting current Go version: %v", err)
			return
		}
		fmt.Printf("Current Go version: %s", currentVersion)
	case "list":
		logger.DebugPrintln("Listing available Go versions...")
		var versionType sources.VersionType
		if len(args) > 1 {
			versionType = sources.VersionType(args[1])
		} else {
			versionType = sources.VersionTypeStable
		}

		var versions []string
		switch versionType {
		case sources.VersionTypeAll:
			logger.DebugPrintln("Listing all versions...")
			versions = sources.GetVersions(sources.VersionTypeAll)
		case sources.VersionTypeLatest:
			logger.DebugPrintln("Listing latest versions...")
			versions = sources.GetVersions(sources.VersionTypeLatest)
		case sources.VersionTypeLts:
			logger.DebugPrintln("Listing LTS versions...")
			versions = sources.GetVersions(sources.VersionTypeLts)
		case sources.VersionTypeStable:
			logger.DebugPrintln("Listing stable versions...")
			versions = sources.GetVersions(sources.VersionTypeStable)
		default:
			fmt.Printf("Unknown version type: %s. Valid types are: latest, lts, stable, all.", versionType)
			return
		}
		fmt.Printf("Available Go versions: %v", versions)
	case "install":
		if len(args) < 2 {
			fmt.Printf("No version specified for install action. Usage: gvm install <version>")
			return
		}
		version := args[1]
		logger.DebugPrintf("Installing Go version: %s", version)
		err := sources.DownloadVersion(version)
		if err != nil {
			fmt.Printf("Error downloading version %s: %v", version, err)
		}
	case "alias":
		if len(args) < 3 {
			fmt.Printf("Not enough arguments for alias action. Usage: gvm alias <source> <target>")
			return
		}
		aliasSource := args[1]
		aliasTarget := args[2]
		if aliasSource == "" || aliasTarget == "" {
			fmt.Printf("Both --alias-source and --alias-target must be specified to create an alias.")
			return
		}
		err := aliases.AddAlias(aliasSource, aliasTarget)
		if err != nil {
			fmt.Printf("Error creating alias from %s to %s: %v", aliasSource, aliasTarget, err)
		}
	case "remove":
		if len(args) < 2 {
			fmt.Printf("No version specified for remove action. Usage: gvm remove <version>")
			return
		}
		version := args[1]
		logger.DebugPrintf("Removing Go version: %s", version)
		err := installed.RemoveVersion(version)
		if err != nil {
			fmt.Printf("Error removing version %s: %v", version, err)
		}
	case "alias-delete":
		if len(args) < 2 {
			fmt.Printf("No alias name specified for alias-delete action. Usage: gvm alias-delete <name>")
			return
		}
		name := args[1]
		logger.DebugPrintf("Removing alias: %s", name)
		err := aliases.DeleteAliasByName(name)
		if err != nil {
			fmt.Printf("Error removing alias %s: %v", name, err)
		}
	case "use":
		if len(args) < 2 {
			fmt.Printf("No version specified for use action. Usage: gvm use <version>")
			return
		}
		version := args[1]
		logger.DebugPrintf("Switching to Go version: %s", version)
		err := aliases.SetVersionAsDefault(version)
		if err != nil {
			fmt.Printf("Error switching to version %s: %v", version, err)
		}
	default:
		fmt.Printf("Unknown action: %s. Valid actions are: list, install, alias.", action)
	}
}

func printHelp() {
	fmt.Println(`gvm: Go Version Manager

Usage:
  gvm <action> [arguments] [flags]

Actions:
  current                  Display the currently active Go version
  list [type]              List available Go versions
                           Types: stable (default), latest, lts, all
  install <version>        Download and install a Go version
  use <version>            Switch to a specific Go version
  alias <source> <target>  Create an alias from one version to another
  alias-delete <name>      Remove an alias by name
  remove <version>         Remove an installed Go version

Flags:
  -h, --help       Show this help message
  -d, --debug      Enable debug output
  -n, --no-cache   Disable caching`)
}
