package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// filter out some files; returns true if the file is a valid resource that should be added to the collection
func isImageFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return ext == ".jpg" && filepath.Base(fileName) != "_.jpg"
}

// if there is no activity on a file for the given timeout, conclude that the file is complete and has been closed
const writeTIMEOUT = 100 * time.Millisecond

// monitor the given folder for new files and add those files as resources (resources.addResource(...))
func StartResGeneratorFromFolder(dir string, webDir string) (err0 error) {
	watcher, err := fsnotify.NewWatcher()
	if err0 = err; err0 != nil {
		return
	}

	err = watcher.Add(dir)
	if err0 = err; err0 != nil {
		return
	}

	go func() {
		defer watcher.Close()
		timeoutChan := make(chan time.Time)
		files := make(map[string]time.Time)

		for {
			select {
			case event, isOpen := <-watcher.Events:
				if !isOpen { return }
				if !isImageFile(event.Name) { continue }

				if false ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Write == fsnotify.Write && !files[event.Name].IsZero() {

					files[event.Name] = time.Now()

					// assume that the file is complete and has been closed if there is no activity for a certain period
					go func() {
						now := <-time.After(writeTIMEOUT)
						timeoutChan <- now
					}()
				}

			case err := <-watcher.Errors:
				log.Fatal("ERROR[PhotoWatcher.err]: ", err)

			case now := <-timeoutChan:
				// if there is no activity for 100ms, we conclude that the file has been closed and it is safe to add it to the collection
				for fileName, lastWriteTime := range files {
					if now.Sub(lastWriteTime) >= writeTIMEOUT {
						delete(files, fileName)
						resources.addResource(fileName, webDir + "/" + filepath.Base(fileName))
					}
				}
			}
		}
	}()

	return
}
