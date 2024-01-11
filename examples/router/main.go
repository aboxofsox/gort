package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.New()

	userData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	router.AddRoute(http.MethodGet, "/", func(c *gort.Context) {
		c.WriteString(http.StatusOK, "Hello World")
	})

	router.AddRoute(http.MethodGet, "/users/:id", func(c *gort.Context) {
		id, ok := c.Params["id"]
		if !ok {
			c.BadRequest()
			return
		}
		user, ok := userData[id]
		if !ok {
			c.NotFound()
			return
		}

		c.WriteString(http.StatusOK, "hello "+user)
	})

	router.AddRoute(http.MethodGet, "/users", func(c *gort.Context) {
		users := make([]string, 0, len(userData))

		for _, user := range userData {
			users = append(users, user)
		}

		if len(users) == 0 {
			c.WriteString(http.StatusNotFound, "no users")
			return
		}

		c.JSON(http.StatusOK, users)
	})

	router.AddRoute(http.MethodGet, "/store/:key", func(c *gort.Context) {
		key, ok := c.Params["key"]
		if !ok {
			c.BadRequest()
			return
		}

		c.Store.Set(key, c.Request().RemoteAddr)
		c.JSON(http.StatusOK, "ok")
	})

	router.AddRoute(http.MethodGet, "/store", func(c *gort.Context) {
		c.JSON(http.StatusOK, c.Store.Items)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))

}
