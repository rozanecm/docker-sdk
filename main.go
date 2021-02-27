package main

import (
	"github.com/jasonlvhit/gocron"
	"sync"
)

const interval = 5

func main() {
	var m sync.Mutex
	leader := ""
	initHttpServer(&leader, &m)
	election(&leader, &m)
	gocron.Start()
	_ = gocron.Every(interval).Second().Do(routineCheck, getNamesOfNodesToControl(), getControlSystemNodeNames(), &leader, &m)
	for {
	}
}
