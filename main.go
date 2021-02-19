package main

import (
	"fmt"
	"os"
)

const thresholdInSeconds = 5

func main() {
	myId := os.Getenv("ID")
	fmt.Printf("My Id: %s\n", myId)
	if iAmLeader(myId) {
		leaderTasks()
	} else {
		noLeaderTasks(myId)
	}
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
}
