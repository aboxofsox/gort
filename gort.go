package gort

import (
	"net/http"
)

type Server struct {
	Router *Router
	Store  *Store
}

func NewServer() *Server {
	return &Server{
		Router: New(),
		Store:  NewStore(),
	}
}

// Handle registers a handler for the given pattern.
func (s *Server) Handle(method, pattern string, handler HandlerFunc) {
	s.Router.AddRoute(method, pattern, handler)
}

// Start starts the server on the given address.
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.Router)
}

// StartTLS starts the server on the given address with TLS.
func (s *Server) StartTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, s.Router)
}

// registerRoute registers a route for the given method and pattern.
func (s *Server) registerRoute(method, pattern string, fileContent []byte, isHTML bool) error {
	s.Router.AddRoute(method, pattern, func(ctx *Context) error {
		if isHTML {
			return ctx.HTML(http.StatusOK, string(fileContent))
		} else {
			return ctx.Send(http.StatusOK, fileContent)
		}
	})
	return nil
}
