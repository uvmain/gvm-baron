package compression

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func unTarGz(source string, target string) error {
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// count total files for progress bar
	tarReader := tar.NewReader(gzReader)
	totalFiles := 0
	for {
		_, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}
		totalFiles++
	}
	file.Close()

	file, err = os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to reopen source file: %w", err)
	}
	defer file.Close()

	gzReader, err = gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader = tar.NewReader(gzReader)

	bar := progressbar.DefaultBytes(
		int64(totalFiles),
		"decompressing",
	)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		bar.Add(1)

		targetPath := filepath.Join(target, header.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(target)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink %s: %w", targetPath, err)
			}
		}
	}

	return nil
}

func unZip(source string, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("error opening zip file: %v", err)
	}
	defer reader.Close()

	bar := progressbar.DefaultBytes(
		int64(len(reader.File)),
		"decompressing",
	)

	for _, file := range reader.File {
		targetPath := fmt.Sprintf("%s%s%s", target, string(os.PathSeparator), file.Name)
		bar.Add(1)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory for file: %v", err)
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening zip entry: %v", err)
		}
		defer rc.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return fmt.Errorf("error copying file content: %v", err)
		}
	}
	return nil
}

func DecompressFile(source string, target string) error {
	extension := filepath.Ext(source)
	switch extension {
	case ".zip":
		return unZip(source, target)
	case ".gz":
		return unTarGz(source, target)
	default:
		return fmt.Errorf("unsupported file extension: %s", extension)
	}
}
