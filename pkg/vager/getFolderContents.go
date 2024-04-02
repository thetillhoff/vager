package vager

import (
	"errors"
	"io/fs"
	"os"
)

func getFolderContents(folderPath string) ([]fs.FileInfo, error) {
	var (
		err error

		folder   *os.File
		fileList []fs.FileInfo
	)

	folder, err = os.Open(folderPath) // Open folder
	if err != nil {
		return fileList, errors.New("failed to open directory: " + err.Error())
	}
	fileList, err = folder.Readdir(0) // Read folder contents
	if err != nil {
		return fileList, errors.New("failed to read directory: " + err.Error())
	}
	folder.Close() // After reading the contents, it can be closed again

	return fileList, nil
}
