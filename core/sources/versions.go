package sources

import (
	"encoding/json"
	"gvm/core/cache"
	"gvm/core/flags"
	"gvm/core/logic"
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

func GetVersions() []string {
	if !flags.NoCache {
		cacheData, cacheAge, err := cache.LoadFromCache(flags.ListType)
		if err == nil {
			if cacheAge < int(time.Hour)*24 {
				logic.DebugPrintf("Using cached versions for type: %s", flags.ListType)
				if items, ok := cacheData.Value.([]interface{}); ok {
					result := make([]string, len(items))
					for i, v := range items {
						result[i] = v.(string)
					}
					return result
				}
			} else {
				logic.DebugPrintf("Cache expired for type: %s, fetching new versions...", flags.ListType)
				cache.DeleteCache(flags.ListType)
			}
		} else {
			logic.DebugPrintf("error fetching cache: %s", err)
		}
	}

	versions := FetchSourceVersions(VersionType(flags.ListType))

	cache.SaveToCache(flags.ListType, versions)
	return versions
}

func FetchSourceVersions(versionType VersionType) []string {
	logic.DebugPrintln("fetching versions from go.dev/dl")
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
		logic.DebugPrintf("error fetching new versions: %s", err)
		return []string{}
	}
	defer response.Body.Close()

	var versions []string
	var data []struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		logic.DebugPrintf("error decoding new versions: %s", err)
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
