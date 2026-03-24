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
var ConfigDirectory string

func InitConfig() {
	Platform = runtime.GOOS
	Arch = runtime.GOARCH
	HomeDirectory, _ = os.UserHomeDir()
	ConfigDirectory = fmt.Sprintf("%s%s%s", HomeDirectory, string(os.PathSeparator), ".gvm-baron")

	_, err := os.Stat(ConfigDirectory)
	if os.IsNotExist(err) {
		log.Println("Config directory does not exist, creating...")
		err := os.Mkdir(ConfigDirectory, 0755)
		if err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}
	}
}
