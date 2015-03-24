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
	flag.StringVar(&dir, "dir", "./", "Path to source tree.")
	flag.Parse()

	var files []string
	var err error

	files, err = findFiles(dir)

	if err != nil {
		log.Fatal(err)
	}

	groups := make([][]string, 0)
	var fileGroups [][]string
	/*
		for _, file := range files {
			fileGroups, err = ParseLineGroups(file)
			if len(fileGroups) > 0 {
				fmt.Println(file)
			}
			for _, fileGroup := range fileGroups {
				groups = append(groups, fileGroup)
			}
		}
	*/
	/*
		for _, group := range groups {
			for _, line := range group {
				fmt.Println(line)
			}
		}
	*/

	fmt.Printf("Run main.")

	//fmt.Printf("%v", groups)
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
		if !isHiddenFile(dir + "/" + file.Name()) {
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

func isHiddenFile(path string) bool {
	for _, part := range strings.Split(path, string(os.PathSeparator)) {
		if part[0:1] == "." && len(part) > 1 {
			return true
		}
	}

	return false
}

/**
 * List Files
 * Get Groupings
 * Parse all Definitions First
 * Parse all Objects
 * Parse all Actions
 * Get Line Type
 * Parse String
 * Parse KeyValue
 * Merge KeyValue
 */
