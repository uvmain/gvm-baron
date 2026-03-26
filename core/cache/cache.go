package cache

import (
	"encoding/json"
	"fmt"
	"gvm/core/config"
	"gvm/core/logger"
	"os"
	"path/filepath"
	"time"
)

type CacheData struct {
	Key       string
	Value     interface{}
	Timestamp int
}

func SaveToCache(key string, value interface{}) error {
	cacheFilePath := getCacheFilePath(key)
	cacheData := CacheData{
		Key:       key,
		Value:     value,
		Timestamp: int(time.Now().Unix()),
	}
	dataBytes, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("error marshaling cache data: %v", err)
	}
	err = os.WriteFile(cacheFilePath, dataBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing cache file: %v", err)
	}
	return nil
}

func DeleteCache(key string) error {
	cacheFilePath := getCacheFilePath(key)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return fmt.Errorf("cache file does not exist")
	}
	err := os.Remove(cacheFilePath)
	if err != nil {
		return fmt.Errorf("error deleting cache file: %v", err)
	}
	return nil
}

func LoadFromCache(key string) (CacheData, int, error) {
	cacheFilePath := getCacheFilePath(key)
	logger.DebugPrintf("loading cache from file %s", cacheFilePath)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return CacheData{}, 0, fmt.Errorf("cache file does not exist")
	}
	file, err := os.Open(cacheFilePath)
	if err != nil {
		return CacheData{}, 0, fmt.Errorf("error opening cache file: %v", err)
	}
	defer file.Close()

	var cacheData CacheData
	if err := json.NewDecoder(file).Decode(&cacheData); err != nil {
		return CacheData{}, 0, fmt.Errorf("error decoding cache data: %v", err)
	}
	age := int(time.Now().Unix()) - cacheData.Timestamp
	logger.DebugPrintf("cache value: %v", cacheData.Value)
	return cacheData, age, nil
}

func getCacheFilePath(key string) string {
	safeKey := fmt.Sprintf("%x", key)
	return filepath.Join(config.CacheDirectory, safeKey+".json")
}
