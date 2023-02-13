package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"log"
	"strconv"
)


// Command to stop the execution of Scaling Manager
var stopCmd = &cobra.Command{
	Use:   "stop",
    Short: "Stop Opensearch Scaling Manager",
    Long:  ``,
	Run: func (cmd *cobra.Command, args []string){
		err:=stop()
		if err!=nil{
			log.Println(err)
			return
		}
		log.Println("Stop Successful")
	},
}

// Input:
// 
// Description:
// 
// 	Function reads the Process Id file and stops the running instance 
//  of Scaling Manager.
// 
// Return:
// 
// (error): Returns error upon unsuccessful execution.
func stop() error{
	log.Println("Stopping Scale Manager")

	_, err := os.Stat("pidFile")
	if err != nil{
		log.Println("Process not found ", err)
		return err
	}

	fileByte, err := os.ReadFile("pidFile")
	if err != nil {
		log.Println("Process id file not found ", err)
		return err
	}

	pid, err := strconv.Atoi(string(fileByte))
	if err != nil {
		log.Println(err)
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Println("Process not found ", err)
		return err
	}

	err = proc.Kill()
	if err != nil {
		log.Println("Unable to terminate process ",err)
		return err
	}

	err = os.Remove("pidFile")
	if err != nil {
		log.Fatalf("Unable to delete pid file ", err)
	}

	return nil
}
