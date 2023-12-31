package vager

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

func Flatten(mainFolderPath string, dryRun bool, verbose bool) {
	var (
		err error
		// fileList fileIO.FileList

		mainFolder   *os.File
		mainFileList []fs.FileInfo

		subFolder   *os.File
		subFileList []fs.FileInfo

		name           string
		resolutionsAsc = []string{
			"480p",
			"720p",
			"1080p",
		}
		highestResolutionsAscIndex int
		cleanedName                string
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

mainFolderLoop:
	for _, mainFolderFile := range mainFileList { // For each file in mainFolder
		name = ""                       // Reset name
		highestResolutionsAscIndex = -1 // Reset highest resolution index

		if mainFolderFile.IsDir() { // If it's a folder
			fullPath := path.Join(mainFolderPath, mainFolderFile.Name())
			if verbose {
				log.Println("Checking folder", fullPath)
			}

			subFolder, err = os.Open(fullPath) // Read subFolder contents
			if err != nil {
				log.Fatalln("failed opening directory", fullPath, err)
			}
			subFileList, err = subFolder.Readdir(0)
			if err != nil {
				log.Fatalln("failed reading directory", fullPath, err)
			}
			subFolder.Close() // After reading the contents, it can be closed again

			if len(subFileList) == 0 { // If this subFolder is empty

				if verbose {
					log.Println("Folder", fullPath, "is empty")
				}

			} else if len(subFileList) == 1 { // If this subFolder only contains one element

				if verbose {
					log.Println("Folder", fullPath, "only contains one element")
				}

				name = subFileList[0].Name()

				cleanedName, err = cleanFileName(name)
				if err != nil {
					log.Fatalln(err)
				}

				// Move file to parent folder
				subFileList[0].Name()
				oldLocation := path.Join(fullPath, name)
				newLocation := path.Join(mainFolderPath, cleanedName)
				if dryRun {
					log.Println("Would move '" + oldLocation + "' to '" + newLocation + "'")
				} else {
					if verbose {
						log.Println("Moving '" + oldLocation + "' to '" + newLocation + "'")
					}
					err := os.Rename(oldLocation, newLocation)
					if err != nil {
						log.Println("Skipping '" + oldLocation + "' because " + errors.Unwrap(err).Error())
						continue // Continue with next folder/file
					}
				}

			} else {
				for _, subFolderfile := range subFileList { // Scan mainFolder contents
					if subFolderfile.IsDir() { // If mainFolder contains a directory
						if verbose {
							log.Println("Skipping folder", fullPath, "because it's contains a folder")
						}
						continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
					}
					if path.Ext(subFolderfile.Name()) != ".mp4" { // If any subFile is not a mp4 file
						if verbose {
							log.Println("Skipping folder", fullPath, "because it contains non-mp4 file", subFolderfile.Name())
						}
						continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
					}

					// TODO Detect resolution from file metadata, not name-suffix

					if name == "" { // If there is no name yet
						for index, resolution := range resolutionsAsc { // Try to match each resolution
							if strings.HasSuffix(subFolderfile.Name(), "_"+resolution+".mp4") { // If resolution matches
								highestResolutionsAscIndex = index                                     // Set this as initial highest
								name = strings.TrimSuffix(subFolderfile.Name(), "_"+resolution+".mp4") // Set this as name to compare other files against
							}
						}
					} else { // If there is a name already
						foundResolution := ""
						for index, resolution := range resolutionsAsc { // Try to match each resolution
							if strings.HasSuffix(subFolderfile.Name(), "_"+resolution+".mp4") { // If resolution matches
								foundResolution = resolution
								if verbose {
									log.Println("Detected resolution", foundResolution)
								}
								if index > highestResolutionsAscIndex { // If found resolution is higher than of the previous files
									highestResolutionsAscIndex = index // Set this as new highest
								}
								if name != strings.TrimSuffix(subFolderfile.Name(), "_"+resolution+".mp4") { // Check name against existing one (from other file/s)
									if verbose {
										log.Println("Skipping folder", fullPath, "because not all of the contained mp4 files have the same name (apart from resolution)")
									}
									continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
								} // Else names match -> ok
							}
						}
						if foundResolution == "" { // If none of the resolutions matched
							if verbose {
								log.Println("Skipping folder", fullPath, "because not all of the contained mp4 files have a resolution set in their file ending (_<resolution>.mp4)")
							}
							continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
						}
					}
					if name == "" { // If there is still no name
						if verbose {
							log.Println("Skipping folder", fullPath, "because none of the contained mp4 files don't have a _<resolution>.mp4 file ending")
						}
						continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
					}

					cleanedName, err = cleanFileName(name)
					if err != nil {
						log.Fatalln(err)
					}

					// TODO Add more cool features :) like
					// - sorting of specific "terms"
					// - reading special attributes like publisher, actors, genre from filename

				}

				// Move file to parent and rename it
				oldLocation := path.Join(fullPath, name+"_"+resolutionsAsc[highestResolutionsAscIndex]+".mp4")
				newLocation := path.Join(mainFolderPath, cleanedName+" - "+resolutionsAsc[highestResolutionsAscIndex]+".mp4")
				if dryRun {
					log.Println("Would move '" + oldLocation + "' to '" + newLocation + "'")
				} else {
					if verbose {
						log.Println("Moving '" + oldLocation + "' to '" + newLocation + "'")
					}
					err := os.Rename(oldLocation, newLocation)
					if err != nil {
						log.Println("Skipping '" + oldLocation + "' because " + errors.Unwrap(err).Error())
						continue // Continue with next folder/file
					} else {
						log.Fatalln(err)
					}
				}
			}

			// TODO Add more cool features :) like
			// - storing special attributes like publisher, actors, genre into metadata of file

			// Delete folder with remaining contents
			if dryRun {
				log.Println("Would delete folder '" + fullPath + "'")
			} else {
				err := os.RemoveAll(fullPath)
				if err != nil {
					log.Fatal(err)
				}
				if verbose {
					log.Println("Deleting folder '" + fullPath + "'")
				}
			}

			log.Println("Flattened '" + fullPath + "'")
			log.Println("---")
		}
	}

	// fileList, err = fileIO.GenerateFileList("./", true)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fileList = fileList.FilterByLevel(0)

	// log.Println(fileList)
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
