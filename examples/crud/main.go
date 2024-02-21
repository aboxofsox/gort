package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	g := gort.New()

	g.GET("/", func(c *gort.Context) error {
		return c.WriteString(200, "hello")
	})

	g.GET("/users", func(c *gort.Context) error {
		return c.JSON(200, c.Store.Items)
	})

	g.GET("/users/:id", func(c *gort.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.BadRequest()
		}

		user, ok := c.Store.Get(id)
		if !ok {
			return c.NotFound()
		}

		return c.WriteString(200, "hello "+user.(string))
	})

	g.GET("/users/:id/delete", func(c *gort.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.BadRequest()
		}

		c.Store.Remove(id)

		return c.WriteString(200, "user deleted")
	})

	g.POST("/users/:id/update", func(c *gort.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.BadRequest()
		}

		name := c.FormValue("name")
		c.Store.Set(id, name)

		return c.WriteString(200, "user updated "+name)
	})

	g.POST("/create", func(c *gort.Context) error {
		c.Request().ParseForm()

		name := c.FormValue("name")
		id := c.FormValue("id")
		if id == "" || name == "" {
			return c.BadRequest()
		}

		c.Store.Set(id, name)

		return c.WriteString(200, "user created: "+name)
	})

	http.ListenAndServe("127.0.0.1:8080", g)
}
