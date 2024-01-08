package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func writeMessage(msg string) gort.HandlerFunc {
	return func(ctx *gort.Context) {
		ctx.WriteString(http.StatusOK, msg)
	}
}

func ping(ctx *gort.Context) {
	ctx.WriteString(http.StatusOK, "pong")
}

func hello(ctx *gort.Context) {
	name := ctx.Params["name"]
	ctx.WriteString(http.StatusOK, "Hello "+name)
}

func handleUser(ctx *gort.Context) {
	userData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	user, ok := userData[ctx.Params["id"]]
	if !ok {
		ctx.WriteString(http.StatusOK, "User not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func main() {
	server := gort.NewServer()

	server.Handle("GET", "/", writeMessage("Hello World"))
	server.Handle("GET", "/ping", ping)
	server.Handle("GET", "/hello/:name", hello)
	server.Handle("GET", "/users/:id", handleUser)

	server.Start(":8080")
}
