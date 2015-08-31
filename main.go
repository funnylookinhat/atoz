package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var PATH_SEPARATOR string = RuneToAscii(os.PathSeparator)

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
	files := make([]string, 0)

	err = filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !isHidden(path) && !file.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func isHidden(path string) bool {
	for i, part := range strings.Split(path, string(PATH_SEPARATOR)) {
		if len(part) > 0 {
			if part[0:1] == "." && i == 0 {
				// Nada
			} else if len(part) > 1 && part[0:2] == ".." {
				// Nada
			} else if part[0:1] == "." && len(part) > 1 {
				return true
			}
		}
	}

	return false
}

// Pretty dang useful - http://stackoverflow.com/a/16684343
func RuneToAscii(r rune) string {
	if r < 128 {
		return string(r)
	} else {
		return "\\u" + strconv.FormatInt(int64(r), 16)
	}
}
