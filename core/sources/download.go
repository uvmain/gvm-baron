package sources

import (
	"fmt"
	"gvm/core/compression"
	"gvm/core/config"
	"gvm/core/files"
	"gvm/core/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	logger.DebugPrintf("Generated file name: %s", fileName)
	return fileName
}

func GenerateDownloadUrl(version string) string {
	fileName := GenerateFileName(version)
	downloadUrl := fmt.Sprintf("https://golang.org/dl/%s", fileName)
	logger.DebugPrintf("Generated download URL: %s", downloadUrl)
	return downloadUrl
}

func DownloadVersion(version string) error {
	downloadUrl := GenerateDownloadUrl(version)
	fileName := GenerateFileName(version)
	logger.DebugPrintf("Downloading version %s from URL: %s", version, downloadUrl)
	tempFilePath := filepath.Join(config.TempDirectory, fileName)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}

	response, err := http.Get(downloadUrl)
	if err != nil {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return fmt.Errorf("failed to download version %s: %w", version, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return fmt.Errorf("failed to download version %s: received non-200 response: %d", version, response.StatusCode)
	}

	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return err
	}

	targetPath := filepath.Join(config.VersionsDirectory, fmt.Sprintf("go%s", version))

	err = compression.DecompressFile(tempFilePath, targetPath)
	if err != nil {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return err
	}

	tempFile.Close()
	err = files.DeleteFile(tempFilePath)
	if err != nil {
		return err
	}

	err = AddAlias(version, "go"+version)
	if err != nil {
		return err
	}

	logger.DebugPrintf("Successfully downloaded and installed Go version %s", version)

	return nil
}

func AddAlias(sourceVersion string, targetVersion string) error {
	if !strings.HasPrefix(sourceVersion, "go") {
		sourceVersion = "go" + sourceVersion
	}
	sourcePath := filepath.Join(config.VersionsDirectory, sourceVersion, "go", "bin", "go")
	targetPath := filepath.Join(config.BinDirectory, targetVersion)

	if files.FileExists(targetPath) {
		logger.DebugPrintf("Alias target already exists: %s. Overwriting alias", targetPath)
		files.DeleteFile(targetPath)
	}

	logger.DebugPrintf("Adding alias for version %s: %s -> %s", sourceVersion, sourcePath, targetPath)
	err := os.Symlink(sourcePath, targetPath)
	if err != nil {
		return err
	}
	return nil
}
