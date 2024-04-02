package vager

// TODO Clean should just clean all names in the specified dir -> Like Flatten, but without looking at subdirs nor moving or deleting files/dirs

import (
	"io/fs"
	"log"
	"os"
)

// Clean iterates over all elements in a folder and clears each name
//
//	clearing means removing invalid characters
//	making all characters lowercase
//	(removing duplicate file-extensions (`.mp4.mp4.mp4`))
func Clean(mainFolderPath string, dryRun bool, verbose bool) {
	var (
		err error

		mainFolder   *os.File
		mainFileList []fs.FileInfo

		cleanedName string
	)

	mainFolder, err = os.Open(mainFolderPath)
	if err != nil {
		log.Fatalln("failed opening directory:", err)
	}
	defer mainFolder.Close()

	mainFileList, err = mainFolder.Readdir(0)
	if err != nil {
		log.Fatalln("failed reading directory:", err)
	}

	for _, mainFolderFile := range mainFileList { // For each file in mainFolder

		if verbose {
			log.Println("Checking name '" + mainFolderFile.Name() + "'")
		}

		cleanedName, err = cleanFileName(mainFolderFile.Name())
		if err != nil {
			log.Fatalln(err)
		}

		// Rename file
		if mainFolderFile.Name() != cleanedName { // If the filename changes
			if dryRun {
				log.Println("Would rename '" + mainFolderFile.Name() + "' to '" + cleanedName + "'")
			} else {
				log.Println("Renaming '" + mainFolderFile.Name() + "' to '" + cleanedName + "'")
				err := os.Rename(mainFolderFile.Name(), cleanedName) // Rename
				if err != nil {
					log.Fatal(err)
				}
			}
		} else { // If the filename stays the same
			if verbose {
				log.Println("Skipping", cleanedName, "as the filename is already clean")
			}
		}

		if verbose {
			log.Println("---")
		}
	}
}
