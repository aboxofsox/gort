package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	g := gort.New()

	g.GET("/", func(c *gort.Context) {
		c.WriteString(200, "hello")
	})

	g.GET("/users", func(c *gort.Context) {
		c.JSON(200, c.Store.Items)
	})

	g.GET("/users/:id", func(c *gort.Context) {
		id := c.Param("id")
		if id == "" {
			c.BadRequest()
			return
		}

		user, ok := c.Store.Get(id)
		if !ok {
			c.NotFound()
			return
		}

		c.WriteString(200, "hello "+user.(string))
	})

	g.GET("/users/:id/delete", func(c *gort.Context) {
		id := c.Param("id")
		if id == "" {
			c.BadRequest()
			return
		}

		c.Store.Remove(id)

		c.WriteString(200, "user deleted")
	})

	g.POST("/users/:id/update", func(c *gort.Context) {
		id := c.Param("id")
		if id == "" {
			c.BadRequest()
			return
		}

		name := c.FormValue("name")
		c.Store.Set(id, name)

		c.WriteString(200, "user updated "+name)
	})

	g.POST("/create", func(c *gort.Context) {
		c.Request().ParseForm()

		name := c.FormValue("name")
		id := c.FormValue("id")
		if id == "" || name == "" {
			c.BadRequest()
			return
		}

		c.Store.Set(id, name)

		c.WriteString(200, "user created: "+name)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", g))
}
