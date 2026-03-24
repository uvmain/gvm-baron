package main

import (
	"gvm/core/config"
	"log"
)

func main() {
	log.Printf("Running on architecture: %s", config.Arch)
	log.Printf("Running on platform: %s", config.Platform)
	log.Printf("User home directory: %s", config.HomeDirectory)
	log.Printf("GVM config directory: %s", config.ConfigDirectory)
}
