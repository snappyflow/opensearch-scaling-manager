package main
  
import (
        "scaling_manager/cmd"
        "os"
        "os/signal"
        "log"
        "syscall"
)

//Input:
// 
//Description:
// 
//	The entry point for the execution of this application
//  The function takes commands(start,stop) to start
//  and stop the Scaling Manager service.
// 
// Return:
func main(){

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c     
        log.Println("Terminating Scale Manager")
        os.Exit(1)
    }()

	err := cmd.Execute()
	if err != nil && err.Error() != "" {
			log.Fatal(err)
	}

}
