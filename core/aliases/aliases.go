package aliases

import (
	"fmt"
	"gvm/core/config"
	"gvm/core/files"
	"gvm/core/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func AddAlias(sourceVersion string, targetName string) error {
	if !strings.HasPrefix(sourceVersion, "go") {
		sourceVersion = "go" + sourceVersion
	}
	sourcePath := filepath.Join(config.VersionsDirectory, sourceVersion, "go", "bin", "go")
	targetPath := filepath.Join(config.BinDirectory, targetName)
	if runtime.GOOS == "windows" {
		sourcePath += ".exe"
		targetPath += ".exe"
	}

	if files.FileExists(targetPath) {
		logger.DebugPrintf("Alias target already exists: %s. Overwriting alias", targetPath)
		files.DeleteFile(targetPath)
	}

	logger.DebugPrintf("Adding alias for version %s: %s -> %s", sourceVersion, sourcePath, targetPath)
	err := os.Symlink(sourcePath, targetPath)
	if err != nil {
		return err
	}
	return nil
}

func AliasForTargetExists(targetVersion string) bool {
	if !strings.HasPrefix(targetVersion, "go") {
		targetVersion = "go" + targetVersion
	}
	targetPath := filepath.Join(config.BinDirectory, targetVersion, "go", "bin", "go")
	targetFiles, err := os.ReadDir(config.BinDirectory)
	if err != nil {
		return false
	}
	for _, file := range targetFiles {
		if file.Type() == os.ModeSymlink {
			linkTarget, err := os.Readlink(filepath.Join(config.BinDirectory, file.Name()))
			if err == nil && linkTarget == targetPath {
				return true
			}
		}
	}
	return false
}

func DeleteAliasesForTarget(targetVersion string) error {
	if !strings.HasPrefix(targetVersion, "go") {
		targetVersion = "go" + targetVersion
	}

	targetPath := filepath.Join(config.BinDirectory, targetVersion, "go", "bin", "go")
	targetFiles, err := os.ReadDir(config.BinDirectory)
	if err != nil {
		return err
	}
	aliasFound := false
	for _, file := range targetFiles {
		if file.Name() == "go" || file.Name() == "go.exe" {
			// don't delete the main 'go' symlink, only version-specific aliases
			continue
		}
		if file.Type() == os.ModeSymlink {
			linkTarget, err := os.Readlink(filepath.Join(config.BinDirectory, file.Name()))
			if err == nil && linkTarget == targetPath {
				aliasFound = true
				err := os.Remove(filepath.Join(config.BinDirectory, file.Name()))
				if err != nil {
					return err
				}
			}
		}
	}
	if !aliasFound {
		return fmt.Errorf("no alias found for target version %s", targetVersion)
	}
	return nil
}

func DeleteAliasByName(aliasName string) error {
	if runtime.GOOS == "windows" && !strings.HasSuffix(aliasName, ".exe") {
		aliasName += ".exe"
	}
	aliasPath := filepath.Join(config.BinDirectory, aliasName)
	_, err := os.Lstat(aliasPath)
	if err != nil {
		return fmt.Errorf("error checking alias %s: %v", aliasPath, err)
	}
	if !os.IsNotExist(err) {
		logger.DebugPrintf("Deleting alias: %s", aliasPath)
		err := files.DeleteFile(aliasPath)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("alias %s does not exist", aliasPath)
	}
}

func SetVersionAsDefault(version string) error {
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}
	versionPath := filepath.Join(config.VersionsDirectory, version, "go", "bin", "go")
	if runtime.GOOS == "windows" {
		versionPath += ".exe"
	}
	if !files.FileExists(versionPath) {
		return fmt.Errorf("version %s is not installed and cannot be set as default", version)
	}

	existingDefaultAlias := filepath.Join(config.BinDirectory, "go")
	if runtime.GOOS == "windows" {
		existingDefaultAlias += ".exe"
	}
	_, err := os.Lstat(existingDefaultAlias)
	if !os.IsNotExist(err) {
		logger.DebugPrintf("Removing existing default link: %s", existingDefaultAlias)
		err := files.DeleteFile(existingDefaultAlias)
		if err != nil {
			return err
		}
		return nil
	}

	err = AddAlias(version, "go")
	if err != nil {
		return fmt.Errorf("error setting version %s as default: %v", version, err)
	}

	logger.DebugPrintf("Version %s set as default: %s -> %s", version, existingDefaultAlias, versionPath)
	return nil
}

func GetCurrentDefaultVersion() (string, error) {
	defaultAlias := filepath.Join(config.BinDirectory, "go")
	if runtime.GOOS == "windows" {
		defaultAlias += ".exe"
	}
	linkTarget, err := os.Readlink(defaultAlias)
	if err != nil {
		return "", fmt.Errorf("error reading default alias: %v", err)
	}
	linkVersion := filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(linkTarget))))
	return linkVersion, nil
}
