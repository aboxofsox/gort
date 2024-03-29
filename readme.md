# gort
~~Another router implementation for Go.~~

Go 1.22 introduces a much better default router, making `gort` more useless than it already was.

## Examples

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

```


## Usage

```go
var Levels = map[Level]string{
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	DEBUG:   "DEBUG",
}
```

#### type Context

```go
type Context struct {
	Params map[string]string
	Writer http.ResponseWriter

	Store  *Store
	Logger *Logger
}
```

Context holdesHTTP request context. It includes parameters, the response writer,
the request, the store, the logger, and a flag indicating whether the response
has been written.

#### func  CreateContext

```go
func CreateContext(w http.ResponseWriter, r *http.Request, store *Store, logger *Logger) *Context
```
CreateContext creates a new context.

#### func (*Context) BadRequest

```go
func (ctx *Context) BadRequest() error
```
BadRequest sets the HTTP status code 400 and writes a message to the response
body.

#### func (*Context) Forbidden

```go
func (ctx *Context) Forbidden() error
```
Forbidden sets the HTTP status code 403 and writes a message to the response
body.

#### func (*Context) FormValue

```go
func (ctx *Context) FormValue(key string) string
```
FormValue returns the value of the given form key.

#### func (*Context) FormValues

```go
func (ctx *Context) FormValues() map[string][]string
```
FormValues returns the values of the form.

#### func (*Context) GetHeader

```go
func (ctx *Context) GetHeader(key string) string
```
GetHeader returns the value of the given header.

#### func (*Context) HTML

```go
func (ctx *Context) HTML(statusCode int, html string) error
```
HTML writes an HTML template the response body.

#### func (*Context) InternalServerError

```go
func (ctx *Context) InternalServerError() error
```
InternalServerError sets the HTTP status code 500 and writes a message to the
response body.

#### func (*Context) JSON

```go
func (ctx *Context) JSON(statusCode int, a any) error
```
SendJSON writes a JSON object to the response body.

#### func (*Context) MethodNotAllowed

```go
func (ctx *Context) MethodNotAllowed() error
```
MethodNotAllowed sets the HTTP status code 405 and writes a message to the
response body.

#### func (*Context) NotFound

```go
func (ctx *Context) NotFound() error
```
NotFound sets the HTTP status code.

#### func (*Context) Param

```go
func (ctx *Context) Param(name string) string
```
Param returns the value of the given parameter.

#### func (*Context) Redirect

```go
func (ctx *Context) Redirect(path string) error
```
Redirect redirects the request to a new URL.

#### func (*Context) Request

```go
func (ctx *Context) Request() *http.Request
```
Request returns the HTTP request.

#### func (*Context) Send

```go
func (ctx *Context) Send(statusCode int, data []byte) error
```
Send writes data to the response body.

#### func (*Context) SetCookie

```go
func (ctx *Context) SetCookie(cookie *http.Cookie)
```
SetCookie sets a cookie in the response.

#### func (*Context) SetHeader

```go
func (ctx *Context) SetHeader(key, value string)
```
SetHeader sets a header in the response.

#### func (*Context) SetHeaders

```go
func (ctx *Context) SetHeaders(headers map[string]string)
```
SetHeaders sets multiple headers in the response.

#### func (*Context) SetStatus

```go
func (ctx *Context) SetStatus(code int)
```
SetStatus sets the HTTP status code.

#### func (*Context) Unauthorized

```go
func (ctx *Context) Unauthorized() error
```
Unauthorized sets the HTTP status code 401 and writes a message to the response
body.

#### func (*Context) WriteString

```go
func (ctx *Context) WriteString(statusCode int, s string) error
```
SendString writes a string to the response body.

#### type Event

```go
type Event struct {
	Timestamp string
	Level     Level
	Message   string
}
```


#### func (Event) String

```go
func (e Event) String() string
```

#### type Group

```go
type Group struct {
}
```


#### func (*Group) AddRoute

```go
func (g *Group) AddRoute(method, pattern string, handler HandlerFunc) error
```

#### func (*Group) DELETE

```go
func (g *Group) DELETE(pattern string, handler HandlerFunc)
```

#### func (*Group) GET

```go
func (g *Group) GET(pattern string, handler HandlerFunc)
```

#### func (*Group) POST

```go
func (g *Group) POST(pattern string, handler HandlerFunc)
```

#### func (*Group) PUT

```go
func (g *Group) PUT(pattern string, handler HandlerFunc)
```

#### type HandlerFunc

```go
type HandlerFunc func(*Context)
```


#### type Level

```go
type Level int
```


```go
const (
	INFO Level = iota
	WARNING
	ERROR
	DEBUG
)
```

#### type Logger

```go
type Logger struct {
}
```


#### func  NewLogger

```go
func NewLogger(w io.Writer) *Logger
```

#### func (*Logger) Debug

```go
func (l *Logger) Debug(message string)
```

#### func (*Logger) Error

```go
func (l *Logger) Error(message string)
```

#### func (*Logger) Info

```go
func (l *Logger) Info(message string)
```

#### func (*Logger) Log

```go
func (l *Logger) Log(level Level, message string)
```

#### func (*Logger) Warning

```go
func (l *Logger) Warning(message string)
```

#### type Route

```go
type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
}
```


#### type Router

```go
type Router struct {
	Logger *Logger
	Groups []*Group
}
```


#### func  New

```go
func New() *Router
```

#### func (*Router) AddRoute

```go
func (r *Router) AddRoute(method, pattern string, handler HandlerFunc)
```
AddRoute adds a new route to the router. It takes the HTTP method, URL pattern,
and handler function as parameters. The method parameter specifies the HTTP
method (e.g., GET, POST, PUT, DELETE). The pattern parameter specifies the URL
pattern that the route should match. The handler parameter is the function that
will be called to handle the request.

#### func (*Router) DELETE

```go
func (r *Router) DELETE(pattern string, handler HandlerFunc)
```

#### func (*Router) Find

```go
func (r *Router) Find(path string) *Route
```
Find returns the route that matches the given path. If no route is found, it
returns nil.

#### func (*Router) GET

```go
func (r *Router) GET(pattern string, handler HandlerFunc)
```

#### func (*Router) Group

```go
func (r *Router) Group(prefix string) *Group
```
Group creates a new group.

#### func (*Router) PATCH

```go
func (r *Router) PATCH(pattern string, handler HandlerFunc)
```

#### func (*Router) POST

```go
func (r *Router) POST(pattern string, handler HandlerFunc)
```

#### func (*Router) PUT

```go
func (r *Router) PUT(pattern string, handler HandlerFunc)
```

#### func(*Router) Static
```go
func (r *Router) Static error
```
Static handles and serves static content. I.E., images, HTML, CSS, JavaScript, etc.

#### func (*Router) ServeHTTP

```go
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
```
ServeHTTP handles the HTTP requests by finding the appropriate route based on
the request URL path, extracting the parameters, and invoking the corresponding
handler. If no route is found, it returns a 404 Not Found response.

#### func (*Router) Use

```go
func (r *Router) Use(handlers ...HandlerFunc)
```
Use adds a new middleware to the router. It takes the middleware function as
parameter. The middleware function is called before the handler function.

#### type Server

```go
type Server struct {
	Router *Router
	Store  *Store
}
```


#### func  NewServer

```go
func NewServer() *Server
```

#### func (*Server) Handle

```go
func (s *Server) Handle(method, pattern string, handler HandlerFunc)
```
Handle registers a handler for the given pattern.

#### func (*Server) Start

```go
func (s *Server) Start(addr string) error
```
Start starts the server on the given address.

#### func (*Server) StartTLS

```go
func (s *Server) StartTLS(addr, certFile, keyFile string) error
```
StartTLS starts the server on the given address with TLS.

#### type Store

```go
type Store struct {
	Items map[string]any
}
```


#### func  NewStore

```go
func NewStore() *Store
```

#### func (*Store) Get

```go
func (s *Store) Get(key string) (any, bool)
```

#### func (*Store) Purge

```go
func (s *Store) Purge()
```

#### func (*Store) Remove

```go
func (s *Store) Remove(key string)
```

#### func (*Store) Set

```go
func (s *Store) Set(key string, value any)
```

#### func (*Store) Size

```go
func (s *Store) Size() int
```
