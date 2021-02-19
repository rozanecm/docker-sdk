package main

import (
	"fmt"
	"os"
)

const thresholdInSeconds = 5

func main() {
	myId := os.Getenv("ID")
	fmt.Printf("My Name: %s\n", myId)
	if iAmLeader(myId) {
		leaderTasks()
	} else {
		nodeName := "node" + myId
		leaderName := "node1"
		noLeaderTasks(nodeName, leaderName)
	}
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
}
