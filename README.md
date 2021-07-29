# IITK Coin

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bhuvansingla/iitk-coin)

IITK Coin is a reward-based pseudo-currency system for the IIT Kanpur campus junta.

## Build and Run

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
