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

	if flags.DebugEnabled {
		logic.DebugPrintln("Debug mode enabled")
	}

	if flags.NoCache {
		logic.DebugPrintln("Cache disabled")
	}

	if flags.ListEnabled {
		logic.DebugPrintln("Listing available Go versions...")
		versions := sources.GetVersions()
		log.Printf("Available Go versions: %v", versions)
	}

	if flags.InstallEnabled {
		if flags.InstallVersion == "" {
			log.Printf("No version specified for installation. Use --install=<version> to specify a version.")
			return
		}
		err := sources.DownloadVersion(flags.InstallVersion)
		if err != nil {
			log.Printf("Error downloading version %s: %v", flags.InstallVersion, err)
		}
	}

	if flags.AliasEnabled {
		if flags.AliasSource == "" || flags.AliasTarget == "" {
			log.Printf("Both --alias-source and --alias-target must be specified to create an alias.")
			return
		}
		err := sources.AddAlias(flags.AliasSource, flags.AliasTarget)
		if err != nil {
			log.Printf("Error creating alias from %s to %s: %v", flags.AliasSource, flags.AliasTarget, err)
		}
	}
}
