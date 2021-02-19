package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func startContainer(name string) {
	fmt.Printf("About to start container: %s\n", name)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}
