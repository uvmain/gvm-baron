package config

import (
	"fmt"
	"gvm/core/logic"
	"os"
	"runtime"
)

var Platform string
var Arch string
var HomeDirectory string
var AppDirectory string
var TempDirectory string
var CacheDirectory string
var VersionsDirectory string
var BinDirectory string
var DebugEnabled bool

func InitConfig() {
	Platform = runtime.GOOS
	Arch = runtime.GOARCH
	HomeDirectory, _ = os.UserHomeDir()
	AppDirectory = fmt.Sprintf("%s%s%s", HomeDirectory, string(os.PathSeparator), ".gvm-baron")

	logic.CreateDir(AppDirectory)

	TempDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "temp")
	logic.CreateDir(TempDirectory)

	CacheDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "cache")
	logic.CreateDir(CacheDirectory)

	VersionsDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "versions")
	logic.CreateDir(VersionsDirectory)

	BinDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "bin")
	logic.CreateDir(BinDirectory)

	logic.DebugPrintf("Initialized config: Platform=%s, Arch=%s, HomeDirectory=%s, AppDirectory=%s, TempDirectory=%s, CacheDirectory=%s, VersionsDirectory=%s, BinDirectory=%s", Platform, Arch, HomeDirectory, AppDirectory, TempDirectory, CacheDirectory, VersionsDirectory, BinDirectory)
}
