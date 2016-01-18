package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilePath represents a file to upload to GCS
type FilePath struct {
	// RelativePath is the path relative to the local root folder from which this FilePath was constructed
	RelativePath string
	// AbsolutePath is the absolute path on disk to the file
	AbsolutePath string
}

func getAllFiles(localRoot string) ([]FilePath, error) {
	if _, err := os.Stat(localRoot); err != nil {
		return err
	}
	var files []FilePath
	err := filepath.Walk(localRoot, func(path string, fInfo os.FileInfo, err error) error {
		if fInfo == nil {
			return fmt.Errorf("got nil FileInfo for path %s", path)
		}
		if fInfo.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(path, localRoot+"/")
		files = append(files, FilePath{AbsolutePath: path, RelativePath: relPath})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
