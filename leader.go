package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"net/http"
	"os"
	"strings"
)

func leaderTasks() {
	nodeNames := getNodesToControlNames()
	//nodes := initNodesInfo(nodeNames)
	gocron.Start()
	_ = gocron.Every(5).Second().Do(routineCheck, nodeNames)
	//initHttpServer(&nodes)
	for {
	}
}

func getNodesToControlNames() []string {
	nodesToControl := os.Getenv("NODES_TO_CONTROL")
	nodesToControl = strings.ReplaceAll(nodesToControl, "\"", "")
	nodesToControlList := strings.Split(nodesToControl, " ")
	fmt.Printf("Nodes to control: %s\n", nodesToControlList)
	return nodesToControlList
}

func routineCheck(nodeNames []string) {
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
