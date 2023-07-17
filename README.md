# video-manager

## Usage

## Installation

## Features

- [ ] iterate through all childfolders in current directory
- [ ] for each folder, check amount of files
  - if there is more than one file
    - if all files are .mp4 files
      - [ ] Check if their names are equal, except for `_<resolution>.mp4`
        - [ ] If that's the case, get the highest resolution, and delete all others. Continue with "if there is only one file"
        - [ ] If that's not the case, skip folder and do nothing
    - if not all files are .mp4 files, skip folder and do nothing
  - if there is only one file
    - [ ] move it one level up -> Rename at the same time, as described in "for all files in current folder that have `.mp4` as file extension"
    - [ ] delete the now empty folder
- [ ] for all files in current folder that have `.mp4` as file extension
  - [ ] remove all "weird" characters
  - [ ] lowercase
  - [ ] search for keywords like names, genre, publisher, sort them, and rename file
  - [ ] add names, genre, publisher to metadata of file
