package main

import (
	"gvm/core/config"
	"gvm/core/flags"
	"gvm/core/logic"
	"gvm/core/sources"
	"log"
)

func main() {
	config.InitConfig()
	flags.InitFlags()

	logic.DebugPrintf("Running on architecture: %s", config.Arch)
	logic.DebugPrintf("Running on platform: %s", config.Platform)
	logic.DebugPrintf("User home directory: %s", config.HomeDirectory)
	logic.DebugPrintf("GVM app directory: %s", config.AppDirectory)

	if flags.ListEnabled {
		logic.DebugPrintf("Fetching available Go versions...")
		versions := sources.GetVersions()
		log.Printf("Available Go versions: %v", versions)
	}

	if flags.DownloadEnabled {
		if flags.DownloadVersion == "" {
			log.Printf("No version specified for download. Use --download=<version> to specify a version.")
			return
		}
		sources.DownloadVersion(flags.DownloadVersion)
	}
}
