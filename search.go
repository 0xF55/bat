package main

import (
	"log"
	"os"
	"strings"
)

// get all bt files in the current directory
func SearchBT() {

	if len(BatFiles) > 0 {
		return
	}

	ext := ".bt"

	dir, err := os.Open("./")
	if err != nil {
		return
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			name := f.Name()
			if strings.HasSuffix(name, ext) {
				BatFiles = append(BatFiles, name)
			}
		}
	}

	// check again

	if len(BatFiles) <= 0 {
		log.Fatal("Can't find bt files")
	}

}
