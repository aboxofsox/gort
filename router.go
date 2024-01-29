package gort

import (
	"net/http"
	"os"
)

type HandlerFunc func(*Context)

type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
}

type Router struct {
	routes      *rtree
	store       *Store
	middlewares []HandlerFunc
	Logger      *Logger
	Groups      []*Group
}

func New() *Router {
	return &Router{
		routes:      newRTree(),
		store:       NewStore(),
		Logger:      NewLogger(os.Stdout),
		middlewares: make([]HandlerFunc, 0),
	}
}

// AddRoute adds a new route to the router.
// It takes the HTTP method, URL pattern, and handler function as parameters.
// The method parameter specifies the HTTP method (e.g., GET, POST, PUT, DELETE).
// The pattern parameter specifies the URL pattern that the route should match.
// The handler parameter is the function that will be called to handle the request.
func (r *Router) AddRoute(method, pattern string, handler HandlerFunc) {
	r.routes.add(&Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
}

// Group creates a new group.
func (r *Router) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
		router: r,
	}

	r.Groups = append(r.Groups, g)

	return g
}

func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.AddRoute(http.MethodGet, pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.AddRoute(http.MethodPost, pattern, handler)
}

func (r *Router) PUT(pattern string, handler HandlerFunc) {
	r.AddRoute(http.MethodPut, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler HandlerFunc) {
	r.AddRoute(http.MethodDelete, pattern, handler)
}

func (r *Router) PATCH(pattern string, handler HandlerFunc) {
	r.AddRoute(http.MethodPatch, pattern, handler)
}

// Use adds a new middleware to the router.
// It takes the middleware function as parameter.
// The middleware function is called before the handler function.
func (r *Router) Use(handlers ...HandlerFunc) {
	r.middlewares = append(r.middlewares, handlers...)
}

// Find returns the route that matches the given path.
// If no route is found, it returns nil.
func (r *Router) Find(path string) *Route {
	return r.routes.find(path)
}

// ServeHTTP handles the HTTP requests by finding the appropriate route based on the request URL path,
// extracting the parameters, and invoking the corresponding handler.
// If no route is found, it returns a 404 Not Found response.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := r.Find(req.URL.Path)
	if route == nil {
		http.NotFound(w, req)
		return
	}

	if route.Method != req.Method {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := &Context{
		Params:  extractParams(req.URL.Path, route.Pattern),
		Writer:  w,
		request: req,
		Store:   r.store,
		Logger:  r.Logger,
	}

	for _, handler := range r.middlewares {
		handler(ctx)
	}

	route.Handler(ctx)
}
