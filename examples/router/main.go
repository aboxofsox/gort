package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.NewRouter()

	userData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	router.AddRoute(http.MethodGet, "/users/:id", func(ctx *gort.Context) {
		id, ok := ctx.Params["id"]
		if !ok {
			ctx.BadRequest()
			return
		}
		user, ok := userData[id]
		if !ok {
			ctx.NotFound()
			return
		}

		ctx.WriteString("hello " + user)
	})

	router.AddRoute(http.MethodGet, "/users", func(ctx *gort.Context) {
		users := make([]string, 0, len(userData))

		for _, user := range userData {
			users = append(users, user)
		}

		if len(users) == 0 {
			ctx.WriteString("no users")
			return
		}

		ctx.JSON(users)
	})

	router.AddRoute(http.MethodGet, "/store/:key", func(ctx *gort.Context) {
		key, ok := ctx.Params["key"]
		if !ok {
			ctx.BadRequest()
			return
		}

		ctx.Store.Set(key, ctx.Request.RemoteAddr)
		ctx.JSON("ok")
	})

	router.AddRoute(http.MethodGet, "/store", func(ctx *gort.Context) {
		ctx.JSON(ctx.Store.Items)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
