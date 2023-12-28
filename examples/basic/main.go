package main

import (
	"github.com/aboxofsox/gort"
)

func writeMessage(msg string) gort.HandlerFunc {
	return func(ctx *gort.Context) {
		ctx.WriteString(msg)
	}
}

func ping(ctx *gort.Context) {
	ctx.WriteString("pong")
}

func hello(ctx *gort.Context) {
	name := ctx.Params["name"]
	ctx.WriteString("Hello " + name)
}

func handleUser(ctx *gort.Context) {
	userData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	user, ok := userData[ctx.Params["id"]]
	if !ok {
		ctx.WriteString("User not found")
		return
	}

	ctx.JSON(user)
}

func main() {
	server := gort.NewServer()

	server.Handle("GET", "/", writeMessage("Hello World"))
	server.Handle("GET", "/ping", ping)
	server.Handle("GET", "/hello/:name", hello)
	server.Handle("GET", "/users/:id", handleUser)

	server.Start(":8080")
}
