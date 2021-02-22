package main

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func leaderTasks() {
	nodeNames := getNodesToControlNames()
	nodes := initNodesInfo(nodeNames)
	gocron.Start()
	_ = gocron.Every(5).Second().Do(routineCheck, nodes)
	initHttpServer(&nodes)
}

func getNodesToControlNames() []string {
	nodesToControl := os.Getenv("NODES_TO_CONTROL")
	nodesToControl = strings.ReplaceAll(nodesToControl, "\"", "")
	nodesToControlList := strings.Split(nodesToControl, " ")
	fmt.Printf("Nodes to control: %s\n", nodesToControlList)
	return nodesToControlList
}

func initNodesInfo(nodeNames []string) map[string]int64 {
	var nodes = make(map[string]int64)
	for _, element := range nodeNames {
		nodes[element] = time.Now().Unix()
	}
	return nodes
}

func routineCheck(nodes map[string]int64) {
	for containerName, timestamp := range nodes {
		fmt.Printf("checking container: %s\n", containerName)
		if timestamp < (time.Now().Unix() - thresholdInSeconds) {
			fmt.Printf("Container %s detected as not running\n", containerName)
			startContainer(containerName)
			fmt.Printf("started container %s\n", containerName)
		}
	}
}

func initHttpServer(nodes *map[string]int64) {
	http.HandleFunc("/heartbeat", heartbeatHandler(nodes))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func heartbeatHandler(nodes *map[string]int64) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		type ExpectedResponse struct {
			Name string `json:"Name"`
		}
		var eR ExpectedResponse
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(request.Body).Decode(&eR)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		nodesCopy := *nodes
		nodesCopy[eR.Name] = time.Now().Unix()
		*nodes = nodesCopy
		fmt.Printf("Nodes after update: %v\n", *nodes)
	}
}
