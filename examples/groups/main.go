package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	g := gort.New()

	apiv1 := g.Group("/api/v1")
	apiv2 := g.Group("/api/v2")
	apiv3 := g.Group("/api/v3")

	apiv1.GET("/users", func(c *gort.Context) {
		c.WriteString(200, "api v1 users")
	})
	apiv2.GET("/users", func(c *gort.Context) {
		c.WriteString(200, "api v2 users")
	})
	apiv3.GET("/users", func(c *gort.Context) {
		c.WriteString(200, "api v3 users")
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", g))

}
