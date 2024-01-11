package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func writeMessage(msg string) gort.HandlerFunc {
	return func(c *gort.Context) {
		c.WriteString(http.StatusOK, msg)
	}
}

func ping(c *gort.Context) {
	c.WriteString(http.StatusOK, "pong")
}

func hello(c *gort.Context) {
	name := c.Params["name"]
	c.WriteString(http.StatusOK, "Hello "+name)
}

func handleUser(c *gort.Context) {
	userData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	user, ok := userData[c.Params["id"]]
	if !ok {
		c.WriteString(http.StatusOK, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func main() {
	server := gort.NewServer()

	server.Handle("GET", "/", writeMessage("Hello World"))
	server.Handle("GET", "/ping", ping)
	server.Handle("GET", "/hello/:name", hello)
	server.Handle("GET", "/users/:id", handleUser)

	server.Start("127.0.0.1:8080")
}
