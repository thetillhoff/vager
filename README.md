# vager
This piece of software aims to clean filestructures.
Primary use is to flatten folders with just one element in them.
In case there are multiple videos in a folder, it will automatically detect the highest resolution, delete the others and flatten the folder.
It also cleans filenames with some rules, like which characters are allowed.

## Usage
```sh
vager flatten
# or
vager clean
```

## Installation
```sh
wget https://github.com/thetillhoff/vager/releases/download/v0.1.0/vager_linux_amd64
wget https://github.com/thetillhoff/vager/releases/download/v0.1.0/vager_linux_amd64.sha256
sha256sum -c vager_linux_amd64.sha256
sudo install vager_linux_amd64 /usr/local/bin/vager # automatically sets rwxr-xr-x permissions
rm vager_linux_amd64 vager_linux_amd64.sha256
```


## Features

- [x] iterate through all childfolders in current directory
- [x] for each folder, check amount of files
  - if there is more than one file
    - if all files are .mp4 files
      - [x] Check if their names are equal, except for `_<resolution>.mp4`
        - [x] If that's the case, get the highest resolution, and delete all others. Continue with "if there is only one file"
        - [x] If that's not the case, skip folder and do nothing
    - if not all files are .mp4 files, skip folder and do nothing
  - if there is only one file
    - [x] move it one level up -> Rename at the same time, as described in "for all files in current folder that have `.mp4` as file extension"
    - [x] delete the now empty folder
- [x] for all files in current folder that have `.mp4` as file extension
  - [x] remove all "weird" characters
  - [x] lowercase
  - [x] search for keywords like names, genre, publisher, sort them, and rename file
  - [x] add names, genre, publisher to metadata of file
