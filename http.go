package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func initHttpServer() {
	http.HandleFunc("/statusCheck", statusCheckHandler)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) } ()
}

func statusCheckHandler(writer http.ResponseWriter, request *http.Request) {
	msg := map[string]string{"Status": "Ok"}
	msgJSON, _ := json.Marshal(msg)
	writer.WriteHeader(200)
	_, _ = writer.Write(msgJSON)
}
