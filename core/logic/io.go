package logic

import (
	"log"
	"os"
)

func CreateDir(directoryPath string) {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		DebugPrintf("%s dir does not exist, creating...", directoryPath)
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

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
