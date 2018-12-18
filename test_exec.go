package main

import (
	"log"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cmd := exec.Command("sleep", "25")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	go func() {
		err = cmd.Wait()
	}()

	var res string
	for {
		if cmd.ProcessState != nil {
			log.Printf(cmd.ProcessState.String())
		} else {
			log.Printf("The process state is nil\n")
		}

		time.Sleep(time.Second * 2)
		if cmd.ProcessState != nil {
			log.Printf(cmd.ProcessState.String())
			status := cmd.ProcessState.Sys().(syscall.WaitStatus)
			switch {
			case status.Exited():
				res = "exit status " + strconv.Itoa(status.ExitStatus())
				log.Printf(res)
			case status.Signaled():
				res = "signal status " + strconv.Itoa(status.ExitStatus())
				log.Printf(res)
			}
		}
		log.Printf("Command finished with error: %v", err)
	}
}
