package sources

import (
	"encoding/json"
	"gvm/core/cache"
	"gvm/core/flags"
	"gvm/core/logic"
	"log"
	"net/http"
	"net/url"
)

type VersionType string

const (
	VersionTypeLatest VersionType = "latest"
	VersionTypeLts    VersionType = "lts"
	VersionTypeStable VersionType = "stable"
	VersionTypeAll    VersionType = "all"
)

func GetVersions() []string {
	packageUrl, err := url.Parse("https://go.dev/dl/")
	if err != nil {
		log.Fatal(err)
	}
	query := packageUrl.Query()
	query.Add("mode", "json")
	if flags.ListType == string(VersionTypeAll) {
		query.Add("include", "all")
	}

	cacheData, cacheAge, err := cache.LoadFromCache(flags.ListType)
	if err == nil {
		if cacheAge < 5*60 {
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
	switch flags.ListType {
	case string(VersionTypeLatest):
		versions = append(versions, data[0].Version)
	case string(VersionTypeLts):
		versions = append(versions, data[1].Version)
	default:
		for _, v := range data {
			versions = append(versions, v.Version)
		}
	}
	cache.SaveToCache(flags.ListType, versions)
	return versions
}
