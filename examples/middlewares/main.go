package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func userMiddleware(users map[string]string) gort.HandlerFunc {
	return func(c *gort.Context) {
		id := c.Param("id")

		user, ok := users[id]
		if !ok {
			c.NotFound()
			return
		}

		c.SetHeader("X-User", user)
	}
}

func loggingMiddleware(c *gort.Context) {
	r := c.Request()
	log.Println(r.Method, r.URL.Path)
}

func main() {
	router := gort.New()

	users := map[string]string{
		"123": "bar",
	}

	router.AddMiddlewares(userMiddleware(users), loggingMiddleware)

	router.AddRoute(http.MethodGet, "/users/:id", func(c *gort.Context) {
		c.WriteString(http.StatusOK, "the user middleware is responsable for setting the X-User header")
	})

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}
