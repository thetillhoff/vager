package videomanager

import (
	"regexp"
	"strings"
)

func cleanFileName(filename string) (string, error) {
	var (
		err error
		reg *regexp.Regexp
	)

	for _, regex := range []string{
		`[^a-zA-Z0-9 \&\#-]+`, // Remove illegal characters in name (by defining allowed ones)
		`[0-9]{4,}`,           // Remove all numbers that are more than 4 characters long
		`\s+`,                 // Remove duplicate spaces
		` ,`,                  // Remove weirdly formatted commas
		`\s+`,                 // Remove duplicate spaces again
		`-$`,                  // Remove trailing hyphen
	} {
		reg, err = regexp.Compile(regex)
		if err != nil {
			return filename, err
		}
		if regex == `\s+` {
			filename = reg.ReplaceAllString(filename, " ") // Replace multiple spaces with single spaces
		} else {
			filename = reg.ReplaceAllString(filename, "") // Remove illegal characters / illegal character combinations
		}
	}

	filename = strings.ToLower(filename) // Make all chars lowercase

	filename = strings.TrimSpace(filename) // Remove leading or trailing whitespace
	return filename, nil
}
