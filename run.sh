#! /bin/bash
docker run -v /var/run/docker.sock:/var/run/docker.sock -it --rm --name my-running-app my-go-app
