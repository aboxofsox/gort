# gort
Another router for Go.

## Install
```bash
go get -u github.com/aboxofsox/gort
```

### Usage
```go
package main

import (
    "net/http"
    "github.com/aboxofsox/gort"
)

func main() {
    // create the server instance
    server := gort.NewServer()
    userStore := map[string]string{
        "foo": "bar",
    }

    server.HandleFunc(http.MethodGet, "/users/:id", func(c *gort.Context){
        id, ok := c.Params["id"]
        if !ok {
            // handle
        }

        user, ok := userStore[id]
        if !ok {
            // handle
        }

        c.JSON(user)
    })

    server.Start(":8080")
}

```

Or just use the router

```go
package main

import (
    "net/http"
    
    "github.com/aboxofsox/gort"
)

func main() {
    router := gort.NewRouter()

    router.AddRoute(http.MethodGet, "/", func(c *gort.Context){
        c.WriteString("Hello World")
    })

    router.AddRoute(http.MethodGet, "/users/:id", func(c *gort.Context){
        id := c.Params["id"]
        if id == "" {
            c.BadRequest()
            return
        }

        c.WriteString(id)
    })

    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Using the Store
The router has a persistent store that can be accessed acrossed contexts.

```go
package main

import (
	"net/http"

	"github.com/aboxofsox/gort"
)

func main() {
	router := gort.NewRouter()

	router.AddRoute(http.MethodGet, "/store/:key/:value", func(c *gort.Context) {
		key, ok := c.Params["key"] // or c.Param("key")
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
		c.JSON("ok")
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

		c.JSON(value)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

```

### Middleware
```go
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
	log.Println(c.Request.Method, c.Request.URL.Path)
}

func main() {
	router := gort.NewRouter()

	users := map[string]string{
		"123": "bar",
	}

    router.AddMiddlewares(userMiddleware(users), loggingMiddleware) // add your middlewares

	router.AddRoute(http.MethodGet, "/users/:id", func(c *gort.Context) {
		c.WriteString("hello world")
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

```

### CRUD
```go
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

    g.POST("/users/:id/update", func(c *gort.Context){
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

```

## Benchmarks
*Copied from [web-framework-benchmark](https://github.com/vishr/web-framework-benchmark)*
```
BenchmarkGortStatic-16             17712             62738 ns/op           39790 B/op        942 allocs/op
BenchmarkGortGitHubAPI-16          10000            114071 ns/op           90369 B/op       1209 allocs/op
BenchmarkGortGplusAPI-16          183368              6456 ns/op            6462 B/op         86 allocs/op
BenchmarkGortParseAPI-16           97270             11991 ns/op            7891 B/op        139 allocs/op
```