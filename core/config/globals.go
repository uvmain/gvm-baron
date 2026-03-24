package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var Platform string
var Arch string
var HomeDirectory string
var AppDirectory string

func InitConfig() {
	Platform = runtime.GOOS
	Arch = runtime.GOARCH
	HomeDirectory, _ = os.UserHomeDir()
	AppDirectory = fmt.Sprintf("%s%s%s", HomeDirectory, string(os.PathSeparator), ".gvm-baron")

	_, err := os.Stat(AppDirectory)
	if os.IsNotExist(err) {
		log.Println("App directory does not exist, creating...")
		err := os.Mkdir(AppDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating app directory: %v", err)
		}
	}
}
