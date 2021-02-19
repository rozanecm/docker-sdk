package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"net/http"
)

func noLeaderTasks(nodeName string, leaderName string) {
	msg := map[string]string{"Name": nodeName}
	msgJSON, _ := json.Marshal(msg)
	gocron.Start()
	_ = gocron.Every(1).Second().Do(sendHeartbeat, msgJSON, leaderName)
	for {
	}
}

func sendHeartbeat(jsonValue []byte, leaderName string) {
	url := "http://" + leaderName + ":8080/heartbeat"
	_, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("an error occurred during POST request: %s\n", err)
	}
}
