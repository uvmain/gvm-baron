package config

import (
	"fmt"
	"os"
	"runtime"
)

var Platform = runtime.GOOS
var Arch = runtime.GOARCH
var HomeDirectory, _ = os.UserHomeDir()
var ConfigDirectory = fmt.Sprintf("%s%s%s", HomeDirectory, string(os.PathSeparator), ".gvm-baron")
