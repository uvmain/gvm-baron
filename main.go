package main

import (
	"gvm/core/config"
	"gvm/core/database"
	"log"
)

func main() {
	config.InitConfig()

	log.Printf("Running on architecture: %s", config.Arch)
	log.Printf("Running on platform: %s", config.Platform)
	log.Printf("User home directory: %s", config.HomeDirectory)
	log.Printf("GVM config directory: %s", config.ConfigDirectory)
	databaseVersion := database.GetDbVersion()
	log.Printf("SQLite version: %s", databaseVersion)
}
