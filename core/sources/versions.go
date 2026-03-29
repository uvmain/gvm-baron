package sources

import (
	"encoding/json"
	"gvm/core/cache"
	"gvm/core/config"
	"gvm/core/logger"
	"log"
	"net/http"
	"net/url"
	"time"
)

type VersionType string

const (
	VersionTypeLatest VersionType = "latest"
	VersionTypeLts    VersionType = "lts"
	VersionTypeStable VersionType = "stable"
	VersionTypeAll    VersionType = "all"
)

type AllVersionTypes struct {
	VersionTypeLatest []string
	VersionTypeLts    []string
	VersionTypeStable []string
	VersionTypeAll    []string
}

func RefreshVersionCache() {
	oldNoCache := config.NoCache
	for _, versionType := range []VersionType{VersionTypeLatest, VersionTypeLts, VersionTypeStable, VersionTypeAll} {
		config.NoCache = true
		cache.DeleteCache(string(versionType))
		logger.DebugPrintf("Refreshing version cache for type: %s", versionType)
		GetVersions(versionType)
	}
	config.NoCache = oldNoCache
}

func GetAllVersionTypes() AllVersionTypes {
	return AllVersionTypes{
		VersionTypeLatest: GetVersions(VersionTypeLatest),
		VersionTypeLts:    GetVersions(VersionTypeLts),
		VersionTypeStable: GetVersions(VersionTypeStable),
		VersionTypeAll:    GetVersions(VersionTypeAll),
	}
}

func GetVersions(versionType VersionType) []string {
	if !config.NoCache {
		cacheData, cacheAge, err := cache.LoadFromCache(string(versionType))
		if err == nil {
			if cacheAge < int(time.Hour)*24 {
				logger.DebugPrintf("Using cached versions for type: %s", versionType)
				if items, ok := cacheData.Value.([]interface{}); ok {
					result := make([]string, len(items))
					for i, v := range items {
						result[i] = v.(string)
					}
					return result
				}
			} else {
				logger.DebugPrintf("Cache expired for type: %s, fetching new versions...", versionType)
				cache.DeleteCache(string(versionType))
			}
		} else {
			logger.DebugPrintf("error fetching cache: %s", err)
		}
	}

	versions := FetchSourceVersions(versionType)

	cache.SaveToCache(string(versionType), versions)
	return versions
}

func FetchSourceVersions(versionType VersionType) []string {
	logger.DebugPrintln("fetching versions from go.dev/dl")
	packageUrl, err := url.Parse("https://go.dev/dl/")
	if err != nil {
		log.Fatal(err)
	}
	query := packageUrl.Query()
	query.Add("mode", "json")
	if versionType == VersionTypeAll {
		query.Add("include", "all")
	}
	packageUrl.RawQuery = query.Encode()
	httpClient := &http.Client{}
	response, err := httpClient.Get(packageUrl.String())
	if err != nil {
		logger.DebugPrintf("error fetching new versions: %s", err)
		return []string{}
	}
	defer response.Body.Close()

	var versions []string
	var data []struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		logger.DebugPrintf("error decoding new versions: %s", err)
		return []string{}
	}
	switch versionType {
	case VersionTypeLatest:
		versions = append(versions, data[0].Version)
	case VersionTypeLts:
		versions = append(versions, data[1].Version)
	default:
		for _, v := range data {
			versions = append(versions, v.Version)
		}
	}
	return versions
}
