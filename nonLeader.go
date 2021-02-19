package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"net/http"
)

func noLeaderTasks(myId string) {
	values := map[string]string{"Id": "node" + myId}
	jsonValue, _ := json.Marshal(values)
	gocron.Start()
	_ = gocron.Every(1).Second().Do(sendHeartbeat, jsonValue)
	for {
	}
}

func sendHeartbeat(jsonValue []byte) {
	_, err := http.Post("http://node1:8080/heartbeat", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("an error occurred during POST request: %s\n", err)
	}
}
