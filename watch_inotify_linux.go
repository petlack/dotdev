//go:build linux
// +build linux

package main

import (
	"log"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// watchFileInotify watches the given file for modifications using inotify,
// and if the file is moved, deleted, or its attributes change, it waits for
// the file to be recreated by watching the parent directory.
func watchFileInotify(filename string, callback func()) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		log.Println("Error initializing inotify:", err)
		return
	}
	defer syscall.Close(fd)

	// Flags to watch for: modifications and events that invalidate the watch.
	flags := uint32(syscall.IN_MODIFY | syscall.IN_MOVE_SELF | syscall.IN_ATTRIB | syscall.IN_DELETE_SELF)
	wd, err := syscall.InotifyAddWatch(fd, filename, flags)
	if err != nil {
		log.Println("Error adding inotify watch on file:", err)
		return
	}

	var buf [4096]byte
	for {
		n, err := syscall.Read(fd, buf[:])
		if err != nil {
			log.Println("Error reading inotify events:", err)
			continue
		}
		if n < syscall.SizeofInotifyEvent {
			continue
		}
		var offset uint32
		for offset <= uint32(n-syscall.SizeofInotifyEvent) {
			event := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
			// On modification, notify clients.
			if event.Mask&syscall.IN_MODIFY != 0 {
				// log.Println("File modified, notifying clients...")
				callback()
			}
			// If the file is moved, its attributes change, or it is deleted, the file watch is no longer valid.
			if event.Mask&(syscall.IN_MOVE_SELF|syscall.IN_ATTRIB|syscall.IN_DELETE_SELF) != 0 {
				// Remove the current file watch.
				syscall.InotifyRmWatch(fd, uint32(wd))

				// Watch the parent directory for creation events.
				parentDir := filepath.Dir(filename)
				dirWd, err := syscall.InotifyAddWatch(fd, parentDir, syscall.IN_CREATE)
				if err != nil {
					log.Println("Error adding inotify watch on directory:", err)
					time.Sleep(500 * time.Millisecond)
					continue
				}
				fileRecreated := false
				// Wait until we see a creation event for our file.
				for !fileRecreated {
					n, err := syscall.Read(fd, buf[:])
					if err != nil {
						log.Println("Error reading inotify events on directory:", err)
						continue
					}
					if n < syscall.SizeofInotifyEvent {
						continue
					}
					var innerOffset uint32
					for innerOffset <= uint32(n-syscall.SizeofInotifyEvent) {
						dirEvent := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[innerOffset]))
						if dirEvent.Mask&syscall.IN_CREATE != 0 {
							// The event's name is stored immediately after the event struct.
							nameBytes := (*[4096]byte)(unsafe.Pointer(&buf[innerOffset+syscall.SizeofInotifyEvent]))[:dirEvent.Len]
							name := strings.TrimRight(string(nameBytes), "\x00")
							if name == filepath.Base(filename) {
								fileRecreated = true
								break
							}
						}
						innerOffset += syscall.SizeofInotifyEvent + dirEvent.Len
					}
				}
				// Remove the directory watch.
				syscall.InotifyRmWatch(fd, uint32(dirWd))
				// Re-add the watch on the file.
				wd, err = syscall.InotifyAddWatch(fd, filename, flags)
				if err != nil {
					log.Println("Error re-adding inotify watch on file:", err)
					time.Sleep(500 * time.Millisecond)
					continue
				}
				callback()
			}
			offset += syscall.SizeofInotifyEvent + event.Len
		}
	}
}
