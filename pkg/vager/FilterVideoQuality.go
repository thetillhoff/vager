package vager

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

// FilterVideoQuality iterates over all elements in a folder and for each
//
//		  checks if it only contains `.mp4` files
//		  if all those files contain a resolution at the end (`_1080p.mp4` or similar)
//	   deletes all but the highest resolution file
func FilterVideoQuality(mainFolderPath string, dryRun bool, verbose bool) {
	var (
		err error

		mainFolder   *os.File
		mainFileList []fs.FileInfo

		subFolder   *os.File
		subFileList []fs.FileInfo

		name           string
		resolutionsAsc = []string{
			"360p",
			"480p",
			"720p",
			"1080p",
		}
		highestResolutionsAscIndex int
		foundResolutions           []string
		cleanedName                string
	)

	mainFolder, err = os.Open(mainFolderPath) // Open main folder
	if err != nil {
		log.Fatalln("failed to open directory:", err)
	}
	mainFileList, err = mainFolder.Readdir(0) // Read main folder contents
	if err != nil {
		log.Fatalln("failed to read directory:", err)
	}
	mainFolder.Close() // After reading the contents, it can be closed again

mainFolderLoop:
	for _, mainFolderElement := range mainFileList { // For each element in mainFolder
		name = ""                       // Reset name
		highestResolutionsAscIndex = -1 // Reset highest resolution index
		foundResolutions = []string{}

		if mainFolderElement.IsDir() { // If it's a folder
			fullPath := path.Join(mainFolderPath, mainFolderElement.Name())
			if verbose {
				log.Println("Checking folder", fullPath)
			}

			subFolder, err = os.Open(fullPath) // Read subFolder contents
			if err != nil {
				log.Fatalln("failed to open directory", fullPath, err)
			}
			subFileList, err = subFolder.Readdir(0)
			if err != nil {
				log.Fatalln("failed to read directory", fullPath, err)
			}
			subFolder.Close() // After reading the contents, it can be closed again

			if len(subFileList) > 1 { // If this subFolder has more than one element
				for _, subFolderFile := range subFileList { // Scan mainFolder contents
					if subFolderFile.IsDir() { // If mainFolder contains a directory
						if verbose {
							log.Println("Skipped folder", fullPath, "because it contains a folder")
							log.Println("---")
						}
						continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
					}
					if path.Ext(subFolderFile.Name()) != ".mp4" { // If any subFile is not a mp4 file
						if verbose {
							log.Println("Skipped folder", fullPath, "because it contains non-mp4 file", subFolderFile.Name())
							log.Println("---")
						}
						continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
					}

					// TODO Detect resolution from file metadata, not name-suffix

					if name == "" { // If there is no name yet
						for index, resolution := range resolutionsAsc { // Try to match each resolution
							if strings.HasSuffix(subFolderFile.Name(), "_"+resolution+".mp4") { // If resolution matches
								highestResolutionsAscIndex = index                                     // Set this as initial highest
								name = strings.TrimSuffix(subFolderFile.Name(), "_"+resolution+".mp4") // Set this as name to compare other files against
							}
						}
					} else { // If there is a name already
						foundResolution := ""
						for index, resolution := range resolutionsAsc { // Try to match each resolution
							if strings.HasSuffix(subFolderFile.Name(), "_"+resolution+".mp4") { // If resolution matches
								foundResolution = resolution
								if verbose {
									log.Println("Detected resolution", foundResolution)
								}
								if index > highestResolutionsAscIndex { // If found resolution is higher than of the previous files
									highestResolutionsAscIndex = index // Set this as new highest
								}
								if name != strings.TrimSuffix(subFolderFile.Name(), "_"+resolution+".mp4") { // Check name against existing one (from other file/s)
									if verbose {
										log.Println("Skipped folder", fullPath, "because not all of the contained mp4 files have the same name (apart from resolution)")
										log.Println("---")
									}
									continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
								} // Else names match -> ok
							}
						}
						if foundResolution == "" { // If none of the resolutions matched
							if verbose {
								log.Println("Skipped folder", fullPath, "because not all of the contained mp4 files have a resolution set in their file ending (_<resolution>.mp4)")
								log.Println("---")
							}
							continue mainFolderLoop // Stop scanning its contents further and continue with next mainFolder (not subFolder!)
						}
						foundResolutions = append(foundResolutions, foundResolution)
					}
					if name == "" { // If there is still no name
						if verbose {
							log.Println("Skipped folder", fullPath, "because none of the contained mp4 files don't have a _<resolution>.mp4 file ending")
							log.Println("---")
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

				// Rename file
				oldLocation := path.Join(fullPath, name+"_"+resolutionsAsc[highestResolutionsAscIndex]+".mp4")
				newLocation := path.Join(fullPath, cleanedName+" - "+resolutionsAsc[highestResolutionsAscIndex]+".mp4")
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
						continue // Continue with next folder/file
					} else {
						log.Fatalln(err)
					}
				}

				// Delete lower-resolution files
				for _, resolution := range foundResolutions[:len(foundResolutions)-1] {
					location := path.Join(fullPath, name+"_"+resolution)
					if dryRun {
						log.Println("Would delete '" + location + "'")
					} else {
						err := os.RemoveAll(fullPath)
						if err != nil {
							log.Fatalln(err)
						}
						if verbose {
							log.Println("Deleted file '" + location + "'")
						}
					}
				}

			} else {
				if verbose {
					log.Println("Skipped folder", fullPath, "because it contains not enough files to filter (<2)")
				}
			}

			// TODO Add more cool features :) like
			// - storing special attributes like publisher, actors, genre into metadata of file

			log.Println("---")
		}
	}
}
