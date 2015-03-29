package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var dir string
	var output string

	flag.StringVar(&dir, "dir", "./", "Path to source tree.")
	flag.StringVar(&output, "output", "./apispec.json", "File to write JSON to.")

	flag.Parse()

	fmt.Println("Reading all files in " + dir)

	var files []string
	var err error

	files, err = findFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found files: ")
	for _, file := range files {
		fmt.Println("\t" + file)
	}

	// apiSpec will be all actions and objects
	// apiSpec := GenerateApiSpec(files)

	// TODO - Write apiSpec out
}

func findFiles(dir string) ([]string, error) {
	var err error
	var dirFiles []os.FileInfo
	var subDirFiles []string

	if dir[len(dir)-1:] == "/" {
		dir = dir[:len(dir)-1]
	}

	dirFiles, err = ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	for _, file := range dirFiles {
		if !isHidden(dir + "/" + file.Name()) {
			if file.IsDir() {
				subDirFiles, err = findFiles(dir + "/" + file.Name())

				if err != nil {
					return nil, err
				}

				if len(subDirFiles) > 0 {
					for _, subDirFile := range subDirFiles {
						files = append(files, subDirFile)
					}
				}
			} else {
				files = append(files, dir+"/"+file.Name())
			}
		}
	}

	return files, nil
}

func isHidden(path string) bool {
	for _, part := range strings.Split(path, string(os.PathSeparator)) {
		if part[0:1] == "." && len(part) > 1 {
			return true
		}
	}

	return false
}
