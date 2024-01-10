package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.New()

	router.AddRoute(http.MethodGet, "/store/:key/:value", func(ctx *gort.Context) {
		key, ok := ctx.Params["key"]
		if !ok {
			ctx.BadRequest()
			return
		}

		value, ok := ctx.Params["value"]
		if !ok {
			ctx.BadRequest()
			return
		}

		ctx.Store.Set(key, value) // set the value in the store
		ctx.JSON(http.StatusOK, "ok")
	})

	router.AddRoute(http.MethodGet, "/store/:key", func(ctx *gort.Context) {
		key, ok := ctx.Params["key"]
		if !ok {
			ctx.BadRequest()
			return
		}

		value, ok := ctx.Store.Get(key) // get the value from the store.
		if !ok {
			ctx.NotFound()
			return
		}

		ctx.JSON(http.StatusOK, value)
	})

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}
