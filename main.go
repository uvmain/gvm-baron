package main

import (
	"gvm/core/config"
	"gvm/core/sources"
	"log"
)

func main() {
	config.InitConfig()

	log.Printf("Running on architecture: %s", config.Arch)
	log.Printf("Running on platform: %s", config.Platform)
	log.Printf("User home directory: %s", config.HomeDirectory)
	log.Printf("GVM app directory: %s", config.AppDirectory)

	versions := sources.GetVersions(sources.VersionTypeLatest)
	log.Printf("Available Go versions: %v", versions)
}
