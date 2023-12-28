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
	routes *rtree
	store  *Store
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

	route.Handler(ctx)
}
