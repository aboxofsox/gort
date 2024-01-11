package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.New()

	router.AddRoute(http.MethodGet, "/store/:key/:value", func(c *gort.Context) {
		key, ok := c.Params["key"]
		if !ok {
			c.BadRequest()
			return
		}

		value, ok := c.Params["value"]
		if !ok {
			c.BadRequest()
			return
		}

		c.Store.Set(key, value) // set the value in the store
		c.JSON(http.StatusOK, "ok")
	})

	router.AddRoute(http.MethodGet, "/store/:key", func(c *gort.Context) {
		key, ok := c.Params["key"]
		if !ok {
			c.BadRequest()
			return
		}

		value, ok := c.Store.Get(key) // get the value from the store.
		if !ok {
			c.NotFound()
			return
		}

		c.JSON(http.StatusOK, value)
	})

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}
