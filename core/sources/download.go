package sources

import (
	"fmt"
	"gvm/core/aliases"
	"gvm/core/compression"
	"gvm/core/config"
	"gvm/core/files"
	"gvm/core/logger"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
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

	headRequest, _ := http.NewRequest("HEAD", downloadUrl, nil)
	headResponse, err := http.DefaultClient.Do(headRequest)
	var contentLength int64 = -1
	if err == nil && headResponse.StatusCode == http.StatusOK {
		contentLength = headResponse.ContentLength
		headResponse.Body.Close()
	}

	request, _ := http.NewRequest("GET", downloadUrl, nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return fmt.Errorf("failed to download version %s: %w", version, err)
	}

	if contentLength == -1 {
		contentLength = response.ContentLength
	}

	var bar *progressbar.ProgressBar
	if contentLength > 0 {
		bar = progressbar.DefaultBytes(contentLength, "downloading")
	} else {
		bar = progressbar.DefaultBytes(-1, "downloading")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		tempFile.Close()
		files.DeleteFile(tempFilePath)
		return fmt.Errorf("failed to download version %s: received non-200 response: %d", version, response.StatusCode)
	}

	_, err = io.Copy(io.MultiWriter(tempFile, bar), response.Body)
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

	err = aliases.AddAlias(version, "go"+version)
	if err != nil {
		return err
	}

	log.Printf("Successfully downloaded and installed Go version %s", version)

	return nil
}
