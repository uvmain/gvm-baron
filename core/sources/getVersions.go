package sources

import (
	"encoding/json"
	"gvm/core/flags"
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
	packageUrl.RawQuery = query.Encode()
	httpClient := &http.Client{}
	response, err := httpClient.Get(packageUrl.String())
	if err != nil {
		return []string{}
	}
	defer response.Body.Close()

	var versions []string
	var data []struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
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
	return versions
}
