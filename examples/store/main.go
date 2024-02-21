package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.New()

	router.AddRoute(http.MethodGet, "/store/:key/:value", func(c *gort.Context) error {
		key, ok := c.Params["key"]
		if !ok {
			return c.BadRequest()

		}

		value, ok := c.Params["value"]
		if !ok {
			return c.BadRequest()

		}

		c.Store.Set(key, value) // set the value in the store
		return c.JSON(http.StatusOK, "ok")
	})

	router.AddRoute(http.MethodGet, "/store/:key", func(c *gort.Context) error {
		key, ok := c.Params["key"]
		if !ok {
			return c.BadRequest()

		}

		value, ok := c.Store.Get(key) // get the value from the store.
		if !ok {
			return c.NotFound()

		}

		return c.JSON(http.StatusOK, value)
	})

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}
