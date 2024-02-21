package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	g := gort.New()
	g.Static("/", "pages")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", g))
}
