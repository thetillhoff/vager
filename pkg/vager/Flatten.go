package vager

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
)

// Flatten iterates over all folders in a folder and
//
//		if it only contains one element, move it one level up and delete the folder
//	    if an element with the same name already exists, add a `-1...` to the filename
//		if it contains no element, delete the folder
//
// //		if it contains multiple files all ending with `.mp4`, moves the highest quality (detected by name suffix `_720p.mp4`) one level up and deletes the folder including the lower quality ones
func Flatten(mainFolderPath string, dryRun bool, verbose bool) {
	var (
		err error

		mainFolderContents []fs.FileInfo

		subFolderPath     string
		subFolderContents []fs.FileInfo

		name        string
		cleanedName string
	)

	mainFolderContents, err = getFolderContents(mainFolderPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, subFolder := range mainFolderContents {
		subFolderPath = path.Join(mainFolderPath, subFolder.Name())
		if subFolder.IsDir() { // Filter for folders
			if verbose {
				log.Println("Checking folder", subFolderPath)
			}

			subFolderContents, err = getFolderContents(subFolderPath)
			if err != nil {
				log.Fatalln(err)
			}

			if len(subFolderContents) == 0 { // If subFolder is empty

				if verbose {
					log.Println("Folder", subFolderPath, "is empty")
				}

				// Delete folder with remaining contents
				if dryRun {
					log.Println("Would delete folder '" + subFolderPath + "'")
				} else {
					err := os.RemoveAll(subFolderPath)
					if err != nil {
						log.Fatalln(err)
					}
					if verbose {
						log.Println("Deleted folder '" + subFolderPath + "'")
					}
				}

			} else if len(subFolderContents) == 1 { // If this subFolder only exactly one element

				if verbose {
					log.Println("Folder", subFolderPath, "only contains exactly one element")
				}

				name = subFolderContents[0].Name()

				cleanedName, err = cleanFileName(name)
				if err != nil {
					log.Fatalln(err)
				}

				// Move file to parent folder
				oldLocation := path.Join(subFolderPath, name)
				newLocation := path.Join(mainFolderPath, cleanedName)
				if dryRun {
					log.Println("Would move '" + oldLocation + "' to '" + newLocation + "'")
				} else {
					if verbose {
						log.Println("Moving '" + oldLocation + "' to '" + newLocation + "'")
					}
					err := os.Rename(oldLocation, newLocation)
					if err != nil {
						log.Println("Skipped '" + oldLocation + "' because " + errors.Unwrap(err).Error())
						log.Println("---")
						continue // Continue with next subfolder
					}
				}

				// Delete folder with remaining contents
				if dryRun {
					log.Println("Would delete folder '" + subFolderPath + "'")
				} else {
					err := os.RemoveAll(subFolderPath)
					if err != nil {
						log.Fatalln(err)
					}
					if verbose {
						log.Println("Deleted folder '" + subFolderPath + "'")
					}
				}

			} else { // If this subFolder contains more than one element
				if verbose {
					log.Println("Skipped folder", subFolderPath, "because it contains more than one element")
				}
			}
		}
	}
}

// TODOs fileIO:
// - GenerateFileList with limited depth (pass amount of levels, ignore all that had more, print them in verbose mode)
// - Get all folders -> filter by type, [both,file,folder]
// - FilterByChildAmount -> type + [==1,>1,<100]
//   -eq 	equals
//   -ne 	not equals
//   -lt 	lower then
//   -le 	lower or equals
//   -gt 	greater then
//   -ge 	greater or equals
// - FilterByExtension
// - FilterByNameRegex
