package gort

import "net/http"

type Group struct {
	prefix string
	router *Router
}

func (g *Group) AddRoute(method, pattern string, handler HandlerFunc) {
	g.router.AddRoute(method, g.prefix+pattern, handler)
}

func (g *Group) GET(pattern string, handler HandlerFunc) {
	g.AddRoute(http.MethodGet, pattern, handler)
}

func (g *Group) POST(pattern string, handler HandlerFunc) {
	g.AddRoute(http.MethodPost, pattern, handler)
}

func (g *Group) PUT(pattern string, handler HandlerFunc) {
	g.AddRoute(http.MethodPut, pattern, handler)
}

func (g *Group) DELETE(pattern string, handler HandlerFunc) {
	g.AddRoute(http.MethodDelete, pattern, handler)
}
