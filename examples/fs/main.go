package main

import (
	"log"

	"github.com/aboxofsox/gort"
)

func main() {
	server := gort.NewServer()
	err := server.FileServer("pages", "/")
	if err != nil {
		log.Fatal(err.Error())
	}
	server.Start("127.0.0.1:8080")
}
