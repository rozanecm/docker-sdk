package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"

	"log"
	"net/http"
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

func noLeaderTasks(myId string) {
	values := map[string]string{"Id": "node" + myId}
	jsonValue, _ := json.Marshal(values)
	gocron.Start()
	_ = gocron.Every(1).Second().Do(sendHeartbeat, jsonValue)
	for {
	}
}

func leaderTasks() {
	nodes := initNodesInfo([]string{"node2","node3","node4"})
	gocron.Start()
	_ = gocron.Every(5).Second().Do(routineCheck, nodes)
	initHttpServer(&nodes)
}

func sendHeartbeat(jsonValue []byte) {
	_, err := http.Post("http://node1:8080/heartbeat", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("an error occurred during POST request: %s\n", err)
	}
}

func initNodesInfo(nodeNames []string) map[string]int64 {
	var nodes = make(map[string]int64)
	for _, element := range nodeNames {
		nodes[element] = time.Now().Unix()
	}
	return nodes
}

func initHttpServer(nodes *map[string]int64) {
	http.HandleFunc("/heartbeat", heartbeatHandler(nodes))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func heartbeatHandler(nodes *map[string]int64) http.HandlerFunc {
	return func(writer http.ResponseWriter, request * http.Request) {
		type ExpectedResponse struct {
			Id string `json:"Id"`
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
		nodesCopy[eR.Id] = time.Now().Unix()
		*nodes = nodesCopy
		fmt.Printf("Nodes after update: %v\n", *nodes)
	}
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
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
