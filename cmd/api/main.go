package main

import (
	"github.com/alexperezortuno/go-url-shortner/cmd/api/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
