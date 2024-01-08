package gort

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type Context struct {
	Params    map[string]string
	Writer    http.ResponseWriter
	Request   *http.Request
	Store     *Store
	Logger    *Logger
	isWritten bool
}

// Param returns the value of the given parameter.
func (ctx *Context) Param(name string) string {
	return ctx.Params[name]
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

// SetCookie sets a cookie in the response.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Writer, cookie)
}

// SetStatus sets the HTTP status code.
func (ctx *Context) SetStatus(code int) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to SetStatus")
		return
	}
	ctx.Writer.WriteHeader(code)
	ctx.isWritten = true
}

// GetHeader returns the value of the given header.
func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

// Send writes data to the response body.
func (ctx *Context) Send(statusCode int, data []byte) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Send")
		return
	}
	ctx.Writer.WriteHeader(statusCode)
	ctx.Writer.Write(data)
	ctx.isWritten = true
}

// SendString writes a string to the response body.
func (ctx *Context) WriteString(statusCode int, s string) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to WriteString")
		return
	}
	ctx.Writer.WriteHeader(statusCode)
	ctx.Writer.Write([]byte(s))
	ctx.isWritten = true
}

// SendJSON writes a JSON object to the response body.
func (ctx *Context) JSON(statusCode int, a any) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to JSON")
		return
	}
	ctx.Writer.Header().Set("Content-Type", "application/json")
	jsn, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Writer.Write(jsn)
	ctx.isWritten = true
}

// SendJSON writes a JSON object to the response body.
func (ctx *Context) HTML(statusCode int, html string) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to HTML")
		return
	}
	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(template.HTML(html)))
	ctx.isWritten = true
}

// Redirect redirects the request to a new URL.
func (ctx *Context) Redirect(path string) {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Redirect")
		return
	}
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusFound)
	ctx.isWritten = true
}

// NotFound sets the HTTP status code.
func (ctx *Context) NotFound() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to NotFound")
		return
	}
	ctx.Writer.WriteHeader(http.StatusNotFound)
	ctx.Writer.Write([]byte("Not Found"))
	ctx.isWritten = true
}

// MethodNotAllowed sets the HTTP status code 405 and writes a message to the response body.
func (ctx *Context) MethodNotAllowed() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to MethodNotAllowed")
		return
	}
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Writer.Write([]byte("Method Not Allowed"))
	ctx.isWritten = true
}

// BadRequest sets the HTTP status code 400 and writes a message to the response body.
func (ctx *Context) BadRequest() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to BadRequest")
		return
	}
	ctx.Writer.WriteHeader(http.StatusBadRequest)
	ctx.Writer.Write([]byte("Bad Request"))
	ctx.isWritten = true
}

// InternalServerError sets the HTTP status code 500 and writes a message to the response body.
func (ctx *Context) InternalServerError() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to InternalServerError")
		return
	}
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write([]byte("Internal Server Error"))
	ctx.isWritten = true
}

// Unauthorized sets the HTTP status code 401 and writes a message to the response body.
func (ctx *Context) Unauthorized() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Unauthorized")
		return
	}
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write([]byte("Unauthorized"))
	ctx.isWritten = true
}

// Forbidden sets the HTTP status code 403 and writes a message to the response body.
func (ctx *Context) Forbidden() {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Forbidden")
		return
	}
	ctx.Writer.WriteHeader(http.StatusForbidden)
	ctx.Writer.Write([]byte("Forbidden"))
	ctx.isWritten = true
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
