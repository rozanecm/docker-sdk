package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func initHttpServer(iAmLeader *bool) {
	http.HandleFunc("/statusCheck", statusCheckHandler)
	http.HandleFunc("/election", electionHandler)
	http.HandleFunc("/leader", leaderHandler(iAmLeader))
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
}

func statusCheckHandler(writer http.ResponseWriter, request *http.Request) {
	msg := map[string]string{"Status": "Ok"}
	msgJSON, _ := json.Marshal(msg)
	writer.WriteHeader(200)
	_, _ = writer.Write(msgJSON)
}

func electionHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("election msg received.")
	//msg := map[string]string{"Status": "Ok"}
	//msgJSON, _ := json.Marshal(msg)
	//writer.WriteHeader(200)
	//_, _ = writer.Write(msgJSON)
}

func leaderHandler(iAmLeader *bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("&&&&&&&&&&&&&&&&&&&&&leader msg received.")
		type ExpectedResponse struct {
			Leader string `json:"Leader"`
		}
		var eR ExpectedResponse
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(request.Body).Decode(&eR)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("**** -> New leader received: %s %s \n", eR.Leader, eR)
		*iAmLeader = false
		if eR.Leader < os.Getenv("NAME") {
			fmt.Println("-> Someone wants to be leader, but I should be!")
			election(nil)
		}
	}
}
