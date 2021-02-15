package main

import (
	"context"
//	"io"
    "fmt"
	"os"

	"github.com/docker/docker/api/types"
//	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
//	"github.com/docker/docker/pkg/stdcopy"

    "github.com/jasonlvhit/gocron"
//    "time"
)

func main() {
    my_id := os.Getenv("ID")
    fmt.Printf("My id: %s\n", my_id)
//    startContainer("ubuntu")
    gocron.Start()
    if (my_id == "1") {
        gocron.Every(5).Second().Do(routineCheck)
    }
//    time.Sleep(30 * time.Second)
    while (true){}

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
    fmt.Printf("%T\n", containers)

    var ids = []string{"/node1", "/node2", "/node3", "/node4"}
    for _, v := range ids {
        fmt.Printf("checking container: %s\n", v)
        if notRunning(v, containers) {
            startContainer(v)
        }
	}
    /*
    for _, v := range containers {
        fmt.Printf("%s\n", v.Names[0])
    }
    */
}

//func notRunning(container string, containers ) {
// do stuff
//}
func notRunning(container string, containers []types.Container) bool {
    for _, currentContainer := range containers {
        if currentContainer.Names[0] == container {
            return false
        }
    }
    return true
}

func startContainer(name string) {
    name = name[1:]
    fmt.Printf("About to start container: %s\n", name)
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }

    /*
     resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: name,
//        Cmd:   []string{"echo", "hello world"},
    }, nil, nil, nil, name)
    if err != nil {
        panic(err)
    }
    */

//    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    if err := cli.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
        panic(err)
    }
 }
