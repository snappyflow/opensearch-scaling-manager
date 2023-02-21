package cmd

import (
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// Input:
//
// Description:
// 	Initializes the stop command, adds the required flags
//
// Return:
func init() {
	stopCmd.PersistentFlags().String("pid", "", "Flag to get the pid")
}

// Command to stop the execution of Scaling Manager
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Opensearch Scaling Manager",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pid, _ := cmd.Flags().GetString("pid")

		if pid != "" {
			err := stop(pid)
			if err != nil {
				log.Error.Println(err)
				return
			}
		} else {
			log.Error.Println("Incorrect Pid")
			return
		}
	},
}

// Input:
//
// Description:
//
//		Function reads the Process Id file and stops the running instance
//	 of Scaling Manager.
//
// Return:
//
// (error): Returns error upon unsuccessful execution.
func stop(pid string) error {
	log.Info.Println("Stopping Scale Manager")
	var pid_int int
	var err error

	pid_int, err = strconv.Atoi(string(pid))
	proc, err := os.FindProcess(pid_int)
	if err != nil {
		log.Error.Println("Process not found ", err)
		return err
	}

	err = proc.Signal(os.Interrupt)
	if err != nil {
		log.Error.Println("Unable to terminate process ", err)
		return err
	}

	time.Sleep(5 * time.Second)

	proc, err = os.FindProcess(pid_int)
	if err != nil {
		log.Info.Println("Process Terminate Successful")
		return nil
	}

	err = proc.Signal(os.Signal(syscall.Signal(0)))
	if err == nil {
		log.Info.Printf("Process with Pid %v is still running.", pid_int)
		log.Info.Println("Scale Manager currently in the provision phase and will be shut down once it is completed")
	} else {
		log.Info.Printf("Process with pid %v is not running.", pid_int)
		log.Info.Println("Process Terminate Successful")
	}
	return nil
}
