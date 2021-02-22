package main

import (
	"fmt"
	"os"
)

const thresholdInSeconds = 5

func main() {
	myId := os.Getenv("ID")
	fmt.Printf("My Name: %s\n", myId)
	initHttpServer()
	if iAmLeader(myId) {
		leaderTasks()
	} else {
		for{}
	}
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
}
