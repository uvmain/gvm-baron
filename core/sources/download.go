package sources

import (
	"fmt"
	"gvm/core/compression"
	"gvm/core/config"
	"gvm/core/logic"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GenerateFileName(version string) string {
	var extension string
	switch config.Platform {
	case "windows":
		extension = "zip"
	default:
		extension = "tar.gz"
	}
	fileName := fmt.Sprintf("go%s.%s-%s.%s", version, config.Platform, config.Arch, extension)
	logic.DebugPrintf("Generated file name: %s", fileName)
	return fileName
}

func GenerateDownloadUrl(version string) string {
	fileName := GenerateFileName(version)
	downloadUrl := fmt.Sprintf("https://golang.org/dl/%s", fileName)
	logic.DebugPrintf("Generated download URL: %s", downloadUrl)
	return downloadUrl
}

func DownloadVersion(version string) error {
	downloadUrl := GenerateDownloadUrl(version)
	fileName := GenerateFileName(version)
	logic.DebugPrintf("Downloading version %s from URL: %s", version, downloadUrl)
	tempFilePath := filepath.Join(config.TempDirectory, fileName)
	out, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(config.VersionsDirectory, fmt.Sprintf("go%s", version))

	err = compression.DecompressFile(tempFilePath, targetPath)
	if err != nil {
		return err
	}

	err = logic.DeleteFile(tempFilePath)
	if err != nil {
		return err
	}

	return nil
}
