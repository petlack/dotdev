package main

import (
	"log"
	"os"
	"time"
)

// watchFile polls the given file for changes and broadcasts a reload message when modified.
func watchFilePoll(filename string, callback func()) {
	var lastModTime time.Time
	for {
		info, err := os.Stat(filename)
		if err != nil {
			log.Println("Error stating file:", err)
		} else {
			modTime := info.ModTime()
			if modTime.After(lastModTime) {
				lastModTime = modTime
				log.Println("File changed, notifying clients...")
				callback()
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
