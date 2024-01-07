package gort

import (
	"net/http"
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
}

func NewRouter() *Router {
	return &Router{
		routes: newRTree(),
		store:  NewStore(),
	}
}

// AddRoute adds a new route to the router.
// It takes the HTTP method, URL pattern, and handler function as parameters.
// The method parameter specifies the HTTP method (e.g., GET, POST, PUT, DELETE).
// The pattern parameter specifies the URL pattern that the route should match.
// The handler parameter is the function that will be called to handle the request.
func (r *Router) AddRoute(method, pattern string, handler HandlerFunc) {
	route := &Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	}
	r.routes.add(route)
}

// AddMiddleware adds a new middleware to the router.
// It takes the middleware function as parameter.
// The middleware function is called before the handler function.
func (r *Router) AddMiddleware(handler HandlerFunc) {
	r.middlewares = append(r.middlewares, handler)
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
		Request: req,
		Store:   r.store,
	}

	for _, handler := range r.middlewares {
		handler(ctx)
	}

	route.Handler(ctx)
}
