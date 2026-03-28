package installed

import (
	"fmt"
	"gvm/core/aliases"
	"gvm/core/config"
	"gvm/core/files"
	"path/filepath"
	"strings"
)

func RemoveVersion(version string) error {
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}

	symlinkExists := aliases.AliasForTargetExists(version)
	if symlinkExists {
		err := aliases.DeleteAliasesForTarget(version)
		if err != nil {
			return err
		}
	}

	installedPath := filepath.Join(config.VersionsDirectory, version)
	if !files.DirectoryExists(installedPath) {
		return fmt.Errorf("version %s is not installed", version)
	}

	err := files.DeleteDirectoryRecursive(installedPath)
	if err != nil {
		return err
	}

	return nil
}
