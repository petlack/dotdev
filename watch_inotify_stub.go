//go:build !linux
// +build !linux

package main

import "log"

// watchFileInotify is a stub for non-Linux platforms.
func watchFileInotify(_ string, _ func()) {
	log.Println("watchFileInotify is not supported on this platform")
}
