package main

import (
	"gvm/core/config"
	"gvm/core/flags"
	"gvm/core/sources"
	"log"
)

func main() {
	config.InitConfig()
	flags.InitFlags()

	log.Printf("Running on architecture: %s", config.Arch)
	log.Printf("Running on platform: %s", config.Platform)
	log.Printf("User home directory: %s", config.HomeDirectory)
	log.Printf("GVM app directory: %s", config.AppDirectory)

	if flags.ListEnabled {
		log.Printf("Fetching available Go versions...")
		versions := sources.GetVersions()
		log.Printf("Available Go versions: %v", versions)
	}
}
