package main

import (
	"bytes"
	"encoding/json"
	"github.com/jasonlvhit/gocron"
	"net/http"
)

const interval = 5

func main() {
	test()
	leader := ""
	initHttpServer(&leader)
	election(&leader)
	gocron.Start()
	_ = gocron.Every(interval).Second().Do(routineCheck, getNamesOfNodesToControl(), getControlSystemNodeNames(), &leader)
	for {
	}
}

func test() {
	/* Test to receive som random env var
	 * to see if they are available when a node is restarted.
	 */
	msg := map[string]string{"First node to control": getControlSystemNodeNames()[0]}
	msgJSON, _ := json.Marshal(msg)

	url := "http://" + "node4" + ":8080/test"
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(msgJSON))
}
