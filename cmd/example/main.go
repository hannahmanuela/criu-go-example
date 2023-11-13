package main

import (
	"fmt"
	"time"
)

func main() {

	// cmd := exec.Command("mount")
	// cmd.Stdout = os.Stdout

	// cmd.Run()

	timer := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-timer.C:
			fmt.Println("exiting")
			return
		default:
			fmt.Println("here sleep")
			time.Sleep(2 * time.Second)
		}
	}

}
