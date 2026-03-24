package main

import (
	"log"
	"runtime"
)

func main() {
	platform := runtime.GOOS
	arch := runtime.GOARCH
	log.Printf("Running on architecture: %s", arch)
	log.Printf("Running on platform: %s", platform)
}
