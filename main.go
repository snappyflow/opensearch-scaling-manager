package main

import (
	"fmt"
	"scaling_manager/cluster"
	"scaling_manager/config"
	log "scaling_manager/logger"
	"scaling_manager/provision"
	"scaling_manager/task"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var state = new(provision.State)

func main() {
	// The following go routine will watch the changes inside config.yaml
	go fileWatch("config.yaml")
	// A periodic check if there is a change in master node to pick up incomplete provisioning
	go periodicProvisionCheck()
	// The polling interval is set to 5 minutes and can be configured.
	ticker := time.Tick(5 * time.Second)
	for range ticker {
		// The recommendation and provisioning should only happen on master node
		if cluster.CheckIfMaster() {
			// This function is responsible for fetching the metrics and pushing it to the index.
			// In starting we will call simulator to provide this details with current timestamp.
			// fetch.FetchMetrics()
			// This function will be responsible for parsing the config file and fill in task_details struct.
			var task = new(task.TaskDetails)
			configStruct, err := config.GetConfig("config.yaml")
			if err != nil {
				log.Error.Println("The recommendation can not be made as there is an error in the validation of config file.")
				log.Error.Println(err.Error())
				continue
			}
			task.Tasks = configStruct.TaskDetails
			// This function is responsible for evaluating the task and recommend.
			recommendationList := task.EvaluateTask()
			// This function is responsible for getting the recommendation and provision.
			provision.GetRecommendation(state, recommendationList)
		}
	}
}

// Input:
// Description: It periodically checks if the master node is changed and picks up if there was any ongoing provision operation
// Output:

func periodicProvisionCheck() {
	tick := time.Tick(5 * time.Second)
	previousMaster := cluster.CheckIfMaster()
	for range tick {
		state.GetCurrentState()
		// Call a function which returns the current master node
		currentMaster := cluster.CheckIfMaster()
		if state.CurrentState != "normal" {
			if !(previousMaster) && currentMaster {
				configStruct, err := config.GetConfig("config.yaml")
				if err != nil {
					log.Warn(log.ProvisionerWarn, "Unable to get Config from GetConfig()")
					return
				}
				cfg := configStruct.ClusterDetails
				if strings.Contains(state.CurrentState, "scaleup") {
					log.Info("Calling scaleOut")
					isScaledUp := provision.ScaleOut(cfg, state)
					if isScaledUp {
						log.Info("Scaleup completed successfully")
					} else {
						// Add a retry mechanism
						log.Warn("Scaleup failed")
					}
				} else if strings.Contains(state.CurrentState, "scaledown") {
					log.Info("Calling scaleIn")
					isScaledDown := provision.ScaleIn(cfg, state)
					if isScaledDown {
						log.Info("Scaledown completed successfully")
					} else {
						// Add a retry mechanism
						log.Warn("Scaledown failed")
					}
				}
			}
		}
		// Update the previousMaster for next loop
		previousMaster = currentMaster
	}
}

func fileWatch(filePath string) {
	//Adding file watcher to detect the change in configuration
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error.Println(fmt.Sprintf("ERROR: %v", err))
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
				log.Error.Println(fmt.Sprintf("EVENT! %#v\n", event))
				log.Error.Println("The recommendation queue will be cleared.")
			case err := <-watcher.Errors:
				log.Error.Println(fmt.Sprintf("ERROR in file watcher: %v", err))
			}
		}
	}()

	// Adding fsnotify watcher to keep track of the changes in config file
	if err := watcher.Add(filePath); err != nil {
		log.Error.Println(fmt.Sprintf("ERROR: %v", err))
	}

	<-done
}
