package installed

import (
	"gvm/core/config"
	"gvm/core/files"
	"strings"
)

func GetInstalledVersions() ([]string, error) {
	entries, err := files.ListDirectory(config.VersionsDirectory)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "go") {
			versions = append(versions, entry.Name())
		}
	}

	return versions, nil
}
