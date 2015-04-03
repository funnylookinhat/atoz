package main

import (
	"encoding/json"
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
	flag.StringVar(&output, "output", "", "File to write JSON to.")

	flag.Parse()

	var files []string
	var err error

	files, err = findFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	var apiSpec ApiSpec
	var resultJson []byte

	apiSpec, err = GenerateApiSpec(files)

	if err != nil {
		log.Fatal(err)
		return
	}

	resultJson, err = json.Marshal(apiSpec)

	if err != nil {
		log.Fatal(err)
		return
	}

	// If no output file specified, throw to stdout
	if len(output) == 0 {
		fmt.Printf("%s", resultJson)
		return
	}

	err = ioutil.WriteFile(output, resultJson, 0644)

	if err != nil {
		log.Fatal(err)
		return
	}
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
