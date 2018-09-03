package main

import (
	"encoding/json"
	"fmt"
	"github.com/amanessinger/manessingercomment-harvester/pkg/harvest"
	"log"
	"os"
)

var l = log.New(os.Stderr, "", log.Lshortfile)

// commandline processing only. Everything else is in pkg
func main() {
	if len(os.Args) != 2 {
		Usage()
	}

	var err error

	baseDir := os.Args[1]
	dirInfo, err := os.Stat(baseDir)
	if os.IsNotExist(err) {
		fmt.Printf("output base directory does not exist")
		Usage()
	}
	if !dirInfo.IsDir() {
		fmt.Printf("output base is no directory")
		Usage()
	}

	var comments []*harvest.Comment
	var receiptHandles []string
	if err, comments, receiptHandles = harvest.FetchComments(); err != nil {
		l.Printf("Error: %v", err)
	}

	for _, c := range comments {
		if err, dirname, filename, createPath := harvest.AttachComment(baseDir, c); err == nil {
			fullDir := fmt.Sprintf("%s/%s", baseDir, dirname)
			if createPath {
				if err := os.MkdirAll(fullDir, 0755); err != nil {
					l.Printf("Can't create directory %s: %v", fullDir, err)
				}
			}
			fullFilename := fmt.Sprintf("%s/%s", fullDir, filename)
			if f, err := os.Create(fullFilename); err != nil {
				l.Printf("Can't create %s: %v", fullFilename, err)
			} else {
				if bytes, err := json.Marshal(*c); err != nil {
					l.Printf("Can't marshal %s: %v", fullFilename, err)
				} else {
					f.Write(bytes)
					l.Printf("%s created\n", fullFilename)
				}
				f.Close()
			}
		} else {
			l.Printf("%v\n", err)
		}
	}

	if len(receiptHandles) > 0 {
		harvest.CleanupQueue(receiptHandles)
	}
}

func Usage() {
	fmt.Printf("Usage: %s <target-base-dir>", os.Args[0])
	os.Exit(1)
}
