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
	nodes := initNodesInfo([]string{"node1","node2","node3","node4"})
	fmt.Printf("My Id: %s\n", myId)
	if iAmLeader(myId) {
		var idsToCheck = getIDsToCheck()
		gocron.Start()
		_ = gocron.Every(5).Second().Do(routineCheck, idsToCheck, nodes)
		initHttpServer(&nodes)
	} else {
		values := map[string]string{"Id": "node" + myId}
		jsonValue, _ := json.Marshal(values)
		gocron.Start()
		_ = gocron.Every(1).Second().Do(sendHeartbeat, jsonValue)
		for {
		}
	}
}

func getIDsToCheck() []string {
	//TODO leer de archivo de conf.
	return []string{"node2", "node3", "node4"}
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
		fmt.Printf("Nodes before: %v\n", *nodes)
		nodesCopy := *nodes
		nodesCopy[eR.Id] = time.Now().Unix()
		*nodes = nodesCopy
		fmt.Printf("Nodes after: %v\n", *nodes)
	}
}

func iAmLeader(id string) bool {
	//TODO hacer esto como corresponde
	return id == "1"
}

func routineCheck(ids []string, nodes map[string]int64) {
	for _, currentID := range ids {
		fmt.Printf("checking container: %s\n", currentID)
		if notRunning(currentID, nodes) {
			fmt.Printf("Container %s detected as not running\n", currentID)
			startContainer(currentID)
			fmt.Printf("started container %s\n", currentID)
		}
	}
}

func notRunning(container string, containers map[string]int64) bool {
	fmt.Printf("Checking node %s. Nodes time: %d; threshold: %d\n",
		container,
		containers[container],
		time.Now().Unix() - thresholdInSeconds)
	return containers[container] < (time.Now().Unix() - thresholdInSeconds)
}
