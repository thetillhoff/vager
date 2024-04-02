package vager

import (
	"io/fs"
	"path"
)

func forEachElementInFolder(folderPath string, fileInfos []fs.FileInfo, dryRun bool, verbose bool, f func(string, fs.FileInfo, bool, bool) error) error {
	var (
		err error
	)

	for _, fileInfo := range fileInfos { // For each element in folder
		err = f(path.Join(folderPath, fileInfo.Name()), fileInfo, dryRun, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}
