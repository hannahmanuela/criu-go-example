package main

import (
	"fmt"
	"time"
)

func main() {

	timer := time.NewTicker(60 * time.Second)

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
