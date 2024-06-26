package vager

import (
	"path"
	"regexp"
	"strings"
)

func cleanFileName(filename string) (string, error) {
	var (
		err error
		ext string
		reg *regexp.Regexp
	)

	ext = path.Ext(filename) // Retrieve file extension (might be empty)

	if ext != "" {
		for strings.HasSuffix(filename, ext) { // If file extension exists at least once
			filename = strings.TrimSuffix(filename, ext) // Remove file extension (it's added later again)
		}
	}

	for _, regex := range []string{
		`[^a-zA-Z0-9 \&\#-,_+]+`, // Remove illegal characters in name (by defining allowed ones)
		`[0-9]{4,}`,              // Remove all numbers that are more than 4 characters long // TODO This causes an error if the filename is _only_ numbers
		`\s+`,                    // Remove duplicate spaces
		` ,`,                     // Remove weirdly formatted commas
		`\s+`,                    // Remove duplicate spaces again
		`-$`,                     // Remove trailing hyphen
		`^-`,                     // Remove beginning hyphen
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

	filename = filename + strings.ToLower(ext) // Append file extension again

	return filename, nil
}
