package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/jasonlvhit/gocron"

	"log"
	"net/http"
)

func main() {
	myId := os.Getenv("ID")
	nodes := initNodesInfo([]string{"node1","node2","node3","node4"})
	fmt.Printf("My Id: %s\n", myId)
	if iAmLeader(myId) {
		gocron.Start()
		_ = gocron.Every(5).Second().Do(routineCheck)
		initHttpServer(&nodes)
	} else {
		values := map[string]string{"Id": "node" + myId}
		jsonValue, _ := json.Marshal(values)
		_, err := http.Post("http://node1:8080/heartbeat", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("an error occurred during POST request: %s\n", err)
		}
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/heartbeat", heartbeatHandler(nodes))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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

func routineCheck() {
	// for name in array with names: if not in docker ps, then start its container.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	//    fmt.Printf("%T\n", containers)

	var ids = []string{"/node1", "/node2", "/node3", "/node4"}
	for _, v := range ids {
		fmt.Printf("checking container: %s\n", v)
		if notRunning(v, containers) {
			startContainer(v)
			fmt.Printf("started container %s\n", v)
		}
	}
	/*
	   for _, v := range containers {
	       fmt.Printf("%s\n", v.Names[0])
	   }
	*/
}

func notRunning(container string, containers []types.Container) bool {
	for _, currentContainer := range containers {
		if currentContainer.Names[0] == container {
			return false
		}
	}
	return true
}
