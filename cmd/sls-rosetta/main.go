package main

import (
	"log"

	"github.com/nikolaymatrosov/sls-rosetta/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
