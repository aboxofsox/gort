package gort

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type Context struct {
	Params  map[string]string
	Writer  http.ResponseWriter
	Request *http.Request
	Store   *Store
}

// SetHeader sets a header in the response.
func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

// SetHeaders sets multiple headers in the response.
func (ctx *Context) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		ctx.SetHeader(k, v)
	}
}

// Send writes data to the response body.
func (ctx *Context) Send(data []byte) {
	ctx.Writer.Write(data)
}

// SendString writes a string to the response body.
func (ctx *Context) WriteString(s string) {
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	ctx.Writer.Write([]byte(s))
}

// SendJSON writes a JSON object to the response body.
func (ctx *Context) JSON(a any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	jsn, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Writer.Write(jsn)
}

// SendJSON writes a JSON object to the response body.
func (ctx *Context) HTML(html string) {
	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.Write([]byte(template.HTML(html)))
}

// Redirect redirects the request to a new URL.
func (ctx *Context) Redirect(path string) {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusFound)
}

// NotFound sets the HTTP status code.
func (ctx *Context) NotFound() {
	ctx.Writer.WriteHeader(http.StatusNotFound)
	ctx.Writer.Write([]byte("Not Found"))
}

// MethodNotAllowed sets the HTTP status code 405 and writes a message to the response body.
func (ctx *Context) MethodNotAllowed() {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Writer.Write([]byte("Method Not Allowed"))
}

// BadRequest sets the HTTP status code 400 and writes a message to the response body.
func (ctx *Context) BadRequest() {
	ctx.Writer.WriteHeader(http.StatusBadRequest)
	ctx.Writer.Write([]byte("Bad Request"))
}

// InternalServerError sets the HTTP status code 500 and writes a message to the response body.
func (ctx *Context) InternalServerError() {
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write([]byte("Internal Server Error"))
}

// Unauthorized sets the HTTP status code 401 and writes a message to the response body.
func (ctx *Context) Unauthorized() {
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write([]byte("Unauthorized"))
}

// Forbidden sets the HTTP status code 403 and writes a message to the response body.
func (ctx *Context) Forbidden() {
	ctx.Writer.WriteHeader(http.StatusForbidden)
	ctx.Writer.Write([]byte("Forbidden"))
}

// extractParams extracts the parameters from the given path based on the provided pattern.
// It takes a path string and a pattern string as input and returns a map[string]string
// containing the extracted parameters. The pattern is expected to have parameter placeholders
// in the form of ":paramName". The function splits the pattern and path into parts and
// matches the corresponding parts to extract the parameters. If the number of parts in the
// pattern and path does not match, it returns nil.
func extractParams(path, pattern string) map[string]string {
	patternParts := strings.Split(pattern, "/")[1:]
	pathParts := strings.Split(path, "/")[1:]

	if len(patternParts) != len(pathParts) {
		return nil
	}

	params := make(map[string]string, len(patternParts))
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			params[strings.TrimPrefix(part, ":")] = pathParts[i]
		}
	}

	return params
}
