package main

import (
	"fmt"
	//	"log"
	"time"
	//	"github.com/influxdata/influxdb/client/v2"
)

type Check struct {
	url string
}

func main() {
	// client := make(map[string]string)

	workCh := make(chan Check, 10)
	stop := make(chan bool)

	go workAssigner(workCh, stop)

	for {
		select {
		case <-time.After(10 * time.Second):
			stop <- true
		}
	}
}

func workAssigner(work chan<- Check, stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("Stopping work queue")
			return
		case <-time.Tick(1 * time.Second):
			fmt.Println("Do stuff here")
			check := Check{"http://localhost:8080/"}
			work <- check
		}
	}
}
