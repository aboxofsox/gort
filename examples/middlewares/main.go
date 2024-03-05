package main

import (
	"log"
	"net/http"

	"github.com/aboxofsox/gort"
)

func userMiddleware(users map[string]string) gort.MiddlewareFunc {
	return func(next gort.HandlerFunc) gort.HandlerFunc {
		return func(c *gort.Context) error {
			id := c.Param("id")

			user, ok := users[id]
			if !ok {
				return c.NotFound()
			}

			c.SetHeader("X-User", user)

			return next(c)
		}
	}
}

func loggingMiddleware(next gort.HandlerFunc) gort.HandlerFunc {
	return func(c *gort.Context) error {
		r := c.Request()
		log.Println(r.Method, r.URL.Path)

		return next(c)
	}
}

func main() {
	router := gort.New()

	users := map[string]string{
		"123": "bar",
	}

	router.Use(userMiddleware(users), loggingMiddleware)

	router.AddRoute(http.MethodGet, "/users/:id", func(c *gort.Context) error {
		return c.WriteString(http.StatusOK, "the user middleware is responsable for setting the X-User header")
	})

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}
