package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const TotalDiskSpace = 70000000
const RequiredDiskSpace = 30000000
const DirectorySizeThreshold = 100000

var rootDirectory *Directory

// struct representing each file
type File struct {
	// file name
	name string
	// file size
	size int
}

// struct representing each directory
type Directory struct {
	// directory name
	name string
	// parent directory pointer
	parent *Directory
	// hash of subdirectory name to subdirectory pointer
	subdirectories map[string]*Directory
	// hash of file name to file pointer
	files map[string]*File
	// cached value of directory size
	// used to avoid computing size multiple times
	// -1 represents empty value
	cachedSize int
}

// directory constructor
func Mkdir() *Directory {
	// create new directory
	directory := &Directory{}
	// initialize hash of subdirectories
	directory.subdirectories = make(map[string]*Directory)
	// initialize hash of files
	directory.files = make(map[string]*File)
	// set cached size to empty value
	directory.cachedSize = -1

	return directory
}

// add file to directory
func (directory *Directory) Touch(file *File) error {
	// check and return error if file already exists
	_, ok := directory.files[file.name]
	if ok {
		return fmt.Errorf("file with name %s already exists in directory %s", file.name, directory.name)
	}

	// store file pointer in hash
	directory.files[file.name] = file

	return nil
}

// add subdirectory to directory
func (directory *Directory) Mv(subdirectory *Directory) error {
	// check and return error if file already exists
	_, ok := directory.subdirectories[subdirectory.name]
	if ok {
		return fmt.Errorf("subdirectory with name %s already exists in directory %s", subdirectory.name, directory.name)
	}

	// store subdirectory pointer in hash
	directory.subdirectories[subdirectory.name] = subdirectory
	// set directory as parent of subdirectory
	subdirectory.parent = directory

	return nil
}

// navigate to directory
func (directory *Directory) Cd(directoryName string) (*Directory, error) {
	// navigate to root directory if absolute root path
	if directoryName == "/" {
		return rootDirectory, nil
	}

	// navigate to parent directory if relative parent path
	if directoryName == ".." {
		if directory.parent == nil {
			return nil, fmt.Errorf("directory %s has no parent", directory.name)
		}

		return directory.parent, nil
	}

	// navigate to subdirectory
	subdirectory, ok := directory.subdirectories[directoryName]
	if !ok {
		return nil, fmt.Errorf("no subdirectory with name %s exists in directory %s", directoryName, directory.name)
	}

	return subdirectory, nil
}

// recursively compute size of directory files and subdirectories
func (directory *Directory) Du() int {
	// if directory cached size is non-empty, return cached size
	if directory.cachedSize >= 0 {
		return directory.cachedSize
	}

	size := 0

	// add file sizes
	for _, file := range directory.files {
		size += file.size
	}

	// recursively add subdirectory sizes
	for _, subdirectory := range directory.subdirectories {
		size += subdirectory.Du()
	}

	// set directory cached size
	directory.cachedSize = size

	return size
}

// get line-by-line file scanner
func GetFileScanner(filePath string) (*bufio.Scanner, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s, %s", filePath, err.Error())
	}

	scanner := bufio.NewScanner(f)

	return scanner, nil
}

func main() {
	// create root directory
	rootDirectory = Mkdir()
	rootDirectory.name = "/"
	// set root directory to current directory
	pwd := rootDirectory
	// list of directories navigated
	directories := make([]*Directory, 0)
	// represents whether or not to expect ls command output in input
	lsMode := false
	// file with filesystem data
	filePath := "../filesystem.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// scan file line-by-line
	for scanner.Scan() {
		input := scanner.Text()
		inputSplit := strings.Fields(input)

		// if input is a command (i.e. cd, ls)
		if inputSplit[0] == "$" {
			cmd := inputSplit[1]

			if cmd == "cd" {
				lsMode = false
				// navigate to directory
				directoryName := inputSplit[2]
				pwd, err = pwd.Cd(directoryName)
			} else if cmd == "ls" {
				// turn on ls mode output
				lsMode = true
			} else {
				panic(fmt.Sprintf("encountered unknown command input %s", input))
			}
		} else if lsMode {
			// if directory listed
			if inputSplit[0] == "dir" {
				directoryName := inputSplit[1]

				// create new directory
				newDirectory := Mkdir()
				newDirectory.name = directoryName

				// add new directory as subdirectory of current directory
				err := pwd.Mv(newDirectory)
				if err != nil {
					panic(fmt.Sprintf("failed to convert add directory %s, %s", directoryName, err.Error()))
				}

				// add new directory to navigated directories
				directories = append(directories, newDirectory)
			} else {
				fileName := inputSplit[1]

				// convert file size from string to int
				fileSize, err := strconv.Atoi(inputSplit[0])
				if err != nil {
					panic(fmt.Sprintf("failed to convert %s into integer, %s", inputSplit[0], err.Error()))
				}

				// create new file
				newFile := &File{
					name: fileName,
					size: fileSize,
				}

				// add new file to current directory
				err = pwd.Touch(newFile)
				if err != nil {
					panic(fmt.Sprintf("failed to convert add directory %s, %s", fileName, err.Error()))
				}
			}
		} else {
			panic(fmt.Sprintf("encountered non-command input %s while lsMode = False", input))
		}
	}

	// compute sizes of navigated directories
	directorySizes := make([]int, len(directories))
	for i := range directories {
		directorySizes[i] = directories[i].Du()
	}

	// get size of directories smaller than threshold
	totalSizeUnderThreshold := 0
	for _, size := range directorySizes {
		if size <= DirectorySizeThreshold {
			totalSizeUnderThreshold += size
		}
	}

	fmt.Println(totalSizeUnderThreshold)

	// compute size of root directory
	rootSize := rootDirectory.Du()
	// compute unused space
	unsusedSpace := TotalDiskSpace - rootSize
	// compute amount of space needed
	extraSpaceNeeded := RequiredDiskSpace - unsusedSpace

	// compute size smallest directory that can be removed to get space needed
	smallestRemovableDirectorySize := rootSize
	for _, size := range directorySizes {
		if size >= extraSpaceNeeded && size < smallestRemovableDirectorySize {
			smallestRemovableDirectorySize = size
		}
	}

	fmt.Println(smallestRemovableDirectorySize)
}
