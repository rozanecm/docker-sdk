package main

import (
	"github.com/jasonlvhit/gocron"
)

const interval = 5

func main() {
	leader := ""
	initHttpServer(&leader)
	election(&leader)
	gocron.Start()
	_ = gocron.Every(interval).Second().Do(routineCheck, getNamesOfNodesToControl(), getControlSystemNodeNames(), &leader)
	for {
	}
}
