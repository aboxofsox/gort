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

    server.HandleFunc(http.MethodGet, "/users/:id", func(ctx *gort.Context){
        id, ok := ctx.Params["id"]
        if !ok {
            // handle
        }

        user, ok := userStore[id]
        if !ok {
            // handle
        }

        ctx.JSON(user)
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

    router.AddRoute(http.MethodGet, "/", func(ctx *gort.Context){
        ctx.WriteString("Hello World")
    })

    router.AddRoute(http.MethodGet, "/users/:id", func(ctx *gort.Context){
        id := ctx.Params["id"]
        if id == "" {
            ctx.BadRequest()
            return
        }

        ctx.WriteString(id)
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
		ctx.JSON("ok")
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

		ctx.JSON(value)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

```