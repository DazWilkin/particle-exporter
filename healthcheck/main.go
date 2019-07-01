// Derived closely from https://github.com/Soluto/golang-docker-healthcheck-example

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected URL as command-line argument")
		os.Exit(1)
	}
	url := os.Args[1]
	if _, err := http.Get(url); err != nil {
		os.Exit(1)
	}
}
