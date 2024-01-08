package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func userMiddleware(users map[string]string) gort.HandlerFunc {
	return func(ctx *gort.Context) {
		id := ctx.Param("id")

		user, ok := users[id]
		if !ok {
			ctx.NotFound()
			return
		}

		ctx.SetHeader("X-User", user)
	}
}

func loggingMiddleware(ctx *gort.Context) {
	log.Println(ctx.Request.Method, ctx.Request.URL.Path)
}

func main() {
	router := gort.New()

	users := map[string]string{
		"123": "bar",
	}

	router.AddMiddlewares(userMiddleware(users), loggingMiddleware)

	router.AddRoute(http.MethodGet, "/users/:id", func(ctx *gort.Context) {
		ctx.WriteString(http.StatusOK, "the user middleware is responsable for setting the X-User header")
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
