package main

import (
	"scaling_manager/cluster"
	"scaling_manager/config"
	"scaling_manager/logger"
	"scaling_manager/task"
	"time"
)

var log logger.LOG

func init() {
	log.Init("logger")
	log.Info.Println("Main module initialized")
}

func main() {
	// The polling interval is set to 5 minutes and can be configured.
	ticker := time.Tick(time.Duration(config.PollingInterval) * time.Second)
	for range ticker {
		// The recommendation and provisioning should only happen on master node.
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
		}
	}
}
