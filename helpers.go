package main

import (
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

func iAmLeader(leader *string) bool {
	return os.Getenv("NAME") == *leader
}
