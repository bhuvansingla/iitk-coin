# IITK Coin

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![Go version](https://img.shields.io/github/go-mod/go-version/bhuvansingla/iitk-coin)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/bhuvansingla/iitk-coin)](https://hub.docker.com/r/bhuvansingla/iitk-coin)
[![Go Build](https://img.shields.io/github/workflow/status/bhuvansingla/iitk-coin/Go?label=go%20build)](https://github.com/bhuvansingla/iitk-coin/actions)
![License](https://img.shields.io/github/license/bhuvansingla/iitk-coin)
![GitHub Repo stars](https://img.shields.io/github/stars/bhuvansingla/iitk-coin)

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

# Update config.yml

# Run it.
./iitk-coin

```

### From Docker

``` bash
# Pull the latest image from DockerHub.
docker pull bhuvansingla/iitk-coin:latest

# Update config.yml

# Run it on port 80 (or whichever you want).
docker run -p 80:8080 -d --name iitk_coin_backend bhuvansingla/iitk-coin

# To access the container shell:
docker exec -t -i iitk_coin_backend /bin/sh

``` 

## Update and Verify User Roles
``` bash
# Ensure that the user, that you want to update the role for, exists.

# Give permissions to execute (if not present already)
chmod +x scripts/*.sh

# Execute the following interactive script to update the role:
bash ./scripts/update-user-role.sh

# To verify user that have a particular role:
bash ./scripts/list-users-by-role.sh

```

## Contributing
The project is open to contributions from the community. For bug fixes and feature requests, please refer to the [CONTRIBUTING](https://github.com/bhuvansingla/iitk-coin/blob/main/.github/CONTRIBUTING.md) guide.

