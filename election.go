package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func election(iAmLeader *bool) {
	iAmNewLeader := startElection(getControlSystemNodeNames())
	if iAmNewLeader {
		fmt.Printf("%d - I, %s, am new leader\n", time.Now().UnixNano(), os.Getenv("NAME"))
		*iAmLeader = true
		announceNewLeadership(getControlSystemNodeNames())
	}
}

func startElection(controlSystemNodeNames []string) bool {
	/* In the election process, one sends an election msg to all
	 * bigger nodes in the system.
	 * If any of them responds, it will continue the election process.
	 * If none of them respond, then the current node is the new leader.
	 * The return statement is the negation of the bool because if no one has answered, then this node is the leader.
	 */
	fmt.Println("In election")
	hasAnyBiggerNodeResponded := false
	for _, node := range controlSystemNodeNames {
		if os.Getenv("NAME") >= node {
			continue
		}
		fmt.Printf("In election with %s\n", node)
		url := "http://" + node + ":8080/election"
		fmt.Printf("Sending election msg to %s\n", node)
		_, err := http.Get(url)
		if err != nil {
			//fmt.Printf("error detected: %s", err)
			fmt.Printf("Container %s detected as not running\n", node)
		} else {
			hasAnyBiggerNodeResponded = true
			/* If a bigger node responded, then there's no need
			 * to send extra messages to other bigger nodes.
			 * This node will continue with the election anyway.
			 */
			break
		}
	}
	return !hasAnyBiggerNodeResponded
}

func announceNewLeadership(controlSystemNodeNames []string) {
	msg := map[string]string{"Leader": os.Getenv("NAME")}
	msgJSON, _ := json.Marshal(msg)
	for _, node := range controlSystemNodeNames {
		if node == os.Getenv("NAME") {
			continue
		}
		url := "http://" + node + ":8080/leader"
		fmt.Printf("Announcing new leader to %s\n", node)
		_, err := http.Post(url, "application/json", bytes.NewBuffer(msgJSON))
		if err != nil {
			fmt.Printf("Error sending leader: %s\n", err)
		}

	}
}
