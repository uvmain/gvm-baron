package main

import (
	"gvm/core/config"
	"log"
)

func main() {
	log.Printf("Running on architecture: %s", config.Arch)
	log.Printf("Running on platform: %s", config.Platform)
}
