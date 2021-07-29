# IITK Coin

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![Go version](https://img.shields.io/github/go-mod/go-version/bhuvansingla/iitk-coin)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/bhuvansingla/iitk-coin)](https://hub.docker.com/r/bhuvansingla/iitk-coin)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/bhuvansingla/iitk-coin/Go/main)](https://github.com/bhuvansingla/iitk-coin/actions)

**IITK Coin** is a reward-based pseudo-currency system for the IIT Kanpur campus junta. 

Detailed vision and regulation rules of this currency are documented in the wiki [here](https://github.com/bhuvansingla/iitk-coin/wiki/Vision-&-Regulation-Rules).

## Build and Run

### From Source
``` bash

# Create the directories if they don't exist already.
cd $GOPATH/src/github.com/bhuvansingla/

# Clone the repository inside.
git clone git@github.com:bhuvansingla/iitk-coin.git
 
cd ./iitk-coin

# Build the project.
go build -o iitk-coin cmd/iitk-coin/main.go

# Run it.
./iitk-coin

```

### From Docker

``` bash
# Pull the latest image from DockerHub.
docker pull bhuvansingla/iitk-coin:latest

# Run it on port 80 (or whichever you want).
docker run -p 80:8080 -d --name iitk_coin_backend bhuvansingla/iitk-coin

# To access the container shell:
docker exec -t -i iitk_coin_backend /bin/sh

```
