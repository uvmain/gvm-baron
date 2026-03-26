package main

import (
	"gvm/core/config"
	"gvm/core/logger"
	"gvm/core/sources"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
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

	if len(args) == 0 {
		log.Println("No action specified. Valid actions are: list, install, alias.")
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

	config.InitConfig()

	action := args[0]

	switch action {
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
			log.Printf("Unknown version type: %s. Valid types are: latest, lts, stable, all.", versionType)
			return
		}
		log.Printf("Available Go versions: %v", versions)
	case "install":
		if len(args) < 2 {
			log.Printf("No version specified for install action. Usage: gvm install <version>")
			return
		}
		version := args[1]
		logger.DebugPrintf("Installing Go version: %s", version)
		err := sources.DownloadVersion(version)
		if err != nil {
			log.Printf("Error downloading version %s: %v", version, err)
		}
	case "alias":
		if len(args) < 3 {
			log.Printf("Not enough arguments for alias action. Usage: gvm alias <source> <target>")
			return
		}
		aliasSource := args[1]
		aliasTarget := args[2]
		if aliasSource == "" || aliasTarget == "" {
			log.Printf("Both --alias-source and --alias-target must be specified to create an alias.")
			return
		}
		err := sources.AddAlias(aliasSource, aliasTarget)
		if err != nil {
			log.Printf("Error creating alias from %s to %s: %v", aliasSource, aliasTarget, err)
		}
	default:
		log.Printf("Unknown action: %s. Valid actions are: list, install, alias.", action)
	}
}
