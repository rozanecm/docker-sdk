package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"os"
	"time"
)

func leaderTasks() {
	nodeNames, err := getNodesToCheckNames("nodes.cfg")
	if err != nil{
		return
	}
	nodes := initNodesInfo(nodeNames)
	gocron.Start()
	_ = gocron.Every(5).Second().Do(routineCheck, nodes)
	initHttpServer(&nodes)
}

func getNodesToCheckNames(path string) ([]string, error) {
	//TODO get this from file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error reading file: %s", err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
	//return []string{"node2", "node3", "node4"}
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
