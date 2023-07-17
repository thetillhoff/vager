package videomanager

// TODO Clean should just clean all names in the specified dir -> Like Flatten, but without looking at subdirs nor moving or deleting files/dirs

import (
	"io/fs"
	"log"
	"os"
)

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

		cleanedName, err = cleanFileName(mainFolderFile.Name())
		if err != nil {
			log.Fatalln(err)
		}

		// Rename file
		if dryRun {
			log.Println("Would move '" + mainFolderFile.Name() + "' to '" + cleanedName + "'")
		} else {
			err := os.Rename(mainFolderFile.Name(), cleanedName)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Cleaned '" + mainFolderFile.Name() + "'")
		log.Println("---")
	}
}
