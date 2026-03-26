package config

import (
	"fmt"
	"gvm/core/files"
	"gvm/core/logger"
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
var NoCache bool

func InitConfig() {
	Platform = runtime.GOOS
	Arch = runtime.GOARCH
	HomeDirectory, _ = os.UserHomeDir()
	AppDirectory = fmt.Sprintf("%s%s%s", HomeDirectory, string(os.PathSeparator), ".gvm-baron")

	files.CreateDir(AppDirectory)

	TempDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "temp")
	files.CreateDir(TempDirectory)

	CacheDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "cache")
	files.CreateDir(CacheDirectory)

	VersionsDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "versions")
	files.CreateDir(VersionsDirectory)

	BinDirectory = fmt.Sprintf("%s%s%s", AppDirectory, string(os.PathSeparator), "bin")
	files.CreateDir(BinDirectory)

	logger.DebugPrintf("Initialized config: Platform=%s, Arch=%s, HomeDirectory=%s, AppDirectory=%s, TempDirectory=%s, CacheDirectory=%s, VersionsDirectory=%s, BinDirectory=%s", Platform, Arch, HomeDirectory, AppDirectory, TempDirectory, CacheDirectory, VersionsDirectory, BinDirectory)
}
