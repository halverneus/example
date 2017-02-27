package main

import (
	"log"

	"github.com/halverneus/example/cli"
)

func main() {
	if err := cli.Parse(); nil != err {
		log.Fatalf("While attempting to parse CLI arguments: %v\n", err)
	}
}
