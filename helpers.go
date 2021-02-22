package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func getNamesOfNodesToControl() []string {
	nodesToControl := os.Getenv("NODES_TO_CONTROL")
	nodesToControl = strings.ReplaceAll(nodesToControl, "\"", "")
	nodesToControlList := strings.Split(nodesToControl, " ")
	fmt.Printf("Nodes to control: %s\n", nodesToControlList)
	return nodesToControlList
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
}

func routineCheck(nodeNames []string, myId string) {
	if !iAmLeader(myId){
		return
	}
	for _, containerName := range nodeNames {
		fmt.Printf("checking container: %s\n", containerName)
		currentURL := "http://" + containerName + ":8080/statusCheck"
		_, err := http.Get(currentURL)
		if err != nil{
			//fmt.Printf("error detected: %s", err)
			fmt.Printf("Container %s detected as not running\n", containerName)
			startContainer(containerName)
			fmt.Printf("started container %s\n", containerName)
		}
	}
}
