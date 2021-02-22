package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
)

const interval = 5

func main() {
	myId := os.Getenv("ID")
	fmt.Printf("My Name: %s\n", myId)
	initHttpServer()
	gocron.Start()
	_ = gocron.Every(interval).Second().Do(routineCheck, getNamesOfNodesToControl(), myId)
	for{}
}
