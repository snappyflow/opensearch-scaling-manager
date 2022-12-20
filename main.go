package main

import (
	"scaling_manager/logger"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// The following go routine will watch the changes inside config.yaml
	go fileWatch("config.yaml")
}

func fileWatch(filePath string) {
	//Adding file watcher to detect the change in configuration
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error.Println("ERROR: ", err)
	}
	defer watcher.Close()
	done := make(chan bool)

	//A go routine that keeps checking for change in configuration
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				//If there is change in config then clear recommendation queue
				//clearRecommendationQueue()
				log.Warn.Println("EVENT!", event)
				log.Warn.Println("The recommendation queue will be cleared.")
			case err := <-watcher.Errors:
				log.Error.Println("ERROR in file watcher: ", err)
			}
		}
	}()

	// Adding fsnotify watcher to keep track of the changes in config file
	if err := watcher.Add(filePath); err != nil {
		log.Error.Println("ERROR:", err)
	}

	<-done
}
