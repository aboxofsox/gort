package gort

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	Router *Router
	Store  *Store
}

func NewServer() *Server {
	return &Server{
		Router: NewRouter(),
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

// FileServer serves files from the given directory.
func (s *Server) FileServer(dir, prefix string) error {
	if _, err := os.Stat(dir); err != nil {
		return err
	}

	readDir(dir, func(path string, entry os.DirEntry) {
		f, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return
		}

		ext := filepath.Ext(path)
		pattern := strings.Replace(filepath.Base(path), ext, "", 1)
		if pattern == "index" {
			pattern = "" // root
		}

		if ext == ".html" || ext == ".htm" {
			s.registerRoute("GET", "/"+prefix+pattern, f, true)
		} else {
			s.registerRoute("GET", "/"+prefix+path, f, false)
		}
	})

	return nil
}

// registerRoute registers a route for the given method and pattern.
func (s *Server) registerRoute(method, pattern string, fileContent []byte, isHTML bool) {
	s.Router.AddRoute(method, pattern, func(ctx *Context) {
		if isHTML {
			ctx.HTML(string(fileContent))
		} else {
			ctx.Send(fileContent)
		}
	})
}
