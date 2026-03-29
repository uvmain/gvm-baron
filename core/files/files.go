package files

import (
	"gvm/core/logger"
	"log"
	"os"
)

func CreateDir(directoryPath string) {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		logger.DebugPrintf("%s dir does not exist, creating...", directoryPath)
		err := os.Mkdir(directoryPath, 0755)
		if err != nil {
			log.Fatalf("Error creating directory %s: %v", directoryPath, err)
		}
	}
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func ListDirectory(directoryPath string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func DirectoryExists(directoryPath string) bool {
	_, err := os.Stat(directoryPath)
	return !os.IsNotExist(err)
}

func DeleteDirectoryRecursive(directoryPath string) error {
	err := os.RemoveAll(directoryPath)
	if err != nil {
		return err
	}
	return nil
}
