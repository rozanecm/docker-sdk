package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

func routineCheck(nodeNames []string, controlSystemNodeNames []string, leader *string, m *sync.Mutex) {
	if iAmLeader(leader, m) {
		m.Lock()
		leaderRoutineCheck(nodeNames)
		m.Unlock()
	} else {
		nonLeaderRoutineCheck(leader, m)
	}
}

func nonLeaderRoutineCheck(leader *string, m *sync.Mutex) {
	currentURL := "http://" + *leader + ":8080/statusCheck"
	_, err := http.Get(currentURL)
	if err != nil {
		//fmt.Printf("error detected: %s", err)
		fmt.Printf("Leader %s detected as not running\n", *leader)
		election(leader, m)
	}
}

func leaderRoutineCheck(nodeNames []string) {
	for _, containerName := range nodeNames {
		if containerName == os.Getenv("NAME") {
			continue
		}
		fmt.Printf("checking container: %s\n", containerName)
		currentURL := "http://" + containerName + ":8080/statusCheck"
		_, err := http.Get(currentURL)
		if err != nil {
			//fmt.Printf("error detected: %s", err)
			fmt.Printf("Container %s detected as not running\n", containerName)
			startContainer(containerName)
			fmt.Printf("started container %s\n", containerName)
		}
	}
}
