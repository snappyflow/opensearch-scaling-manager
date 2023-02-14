package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"os/exec"
	"log"
	app "github.com/maplelabs/opensearch-scaling-manager/scaleManager"
)

// Start Command to start the Scaling Manager service
var startCmd = &cobra.Command{
	Use:   "start",
    Short: "Start Opensearch Scaling Manager",
    Long:  ``,
	Run: func (cmd *cobra.Command, args []string){
		bg,_:=cmd.Flags().GetString("b")

		if bg!=""{
				if bg != "true"{
					log.Fatal("Unknown value for flag")
					os.Exit(1)
				}
				log.Println("Started scale manager in background")
				startBackground()
		}else{
			err:=start()
			if err!=nil{
				log.Fatal(err)
			}
		}
	},
}

// Input:
// 
// Description:
// 
// 	The Function initilazes and starts the execution of Scaling Manager
// 
// Return:
// 
// 	(error): Returns error upon unsuccessful execution
func start() error{
	app.Initialize()
	app.Run()
	return nil
}

// Input:
// 
// Description:
// 
// 	The Function is executed when user sets flag(--b=true) along
// 	with start command. It creates a Background process and creates
// 	a file to track the Process Id of the background process. 
// 
// Return:
// 
// 	(error): Returns error upon unsuccessful execution.
func startBackground() error{
	_, err := os.Stat("pidFile")

	if err != nil{
		scaleManagerExe, err := os.Executable()
		if err != nil{
			log.Println(err)
			return err
		}

		cmd := exec.Command(scaleManagerExe, "start")

		err = cmd.Start()
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("Scale Manager started with pid %v", cmd.Process.Pid)

		err = os.WriteFile("pidFile", []byte(fmt.Sprintf("%v", cmd.Process.Pid)), 0644)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	log.Println("Process already running")
	fileByte, err := os.ReadFile("pidFile")
	if err != nil {
		log.Println(err)
	}

	log.Println("Process already running with pid %v ", string(fileByte))
	return nil
}

// Input:
// 
// Description:
// 	Initializes the start command, adds the required flags
// 
// Return:
func init(){
	startCmd.PersistentFlags().String("b","","Flag to run process in background")
}
