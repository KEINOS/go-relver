package main

import (
	"fmt"
	"log"

	"github.com/KEINOS/go-relver/relver"
)

func main() {
	// Get the latest release version of Go
	goVerLatest, err := relver.Get()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Current latest go version is: " + goVerLatest)
}
