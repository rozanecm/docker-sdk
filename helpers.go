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
	//fmt.Printf("Nodes to control: %s\n", nodesToControlList)
	return nodesToControlList
}

func getControlSystemNodeNames() []string {
	controlSystemNodes := os.Getenv("CONTROL_SYSTEM_NODES")
	controlSystemNodes = strings.ReplaceAll(controlSystemNodes, "\"", "")
	controlSystemNodesList := strings.Split(controlSystemNodes, " ")
	//fmt.Printf("Control system Nodes: %s\n", controlSystemNodesList)
	return controlSystemNodesList
}

func routineCheck(nodeNames []string, controlSystemNodeNames []string, iAmLeader bool) {
	if !iAmLeader {
		return
	}
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
