package gort

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Context holdesHTTP request context. It includes parameters,
// the response writer, the request, the store, the logger, and a flag
// indicating whether the response has been written.
type Context struct {
	Params    map[string]string
	Writer    http.ResponseWriter
	request   *http.Request
	Store     *Store
	Logger    *Logger
	isWritten bool
}

// CreateContext creates a new context.
func CreateContext(w http.ResponseWriter, r *http.Request, store *Store, logger *Logger) *Context {
	return &Context{
		Params:  make(map[string]string),
		Writer:  w,
		request: r,
		Store:   store,
		Logger:  logger,
	}
}

// setParams sets the parameters for the context.
func (ctx *Context) setParams(params map[string]string) {
	ctx.Params = params
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
	return ctx.request.Header.Get(key)
}

// Send writes data to the response body.
func (ctx *Context) Send(statusCode int, data []byte) error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Send")
		return errors.New("superflous call to Send")
	}
	_, err := write(ctx, statusCode, data)
	if err != nil {
		return fmt.Errorf("error writing data to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// SendString writes a string to the response body.
func (ctx *Context) WriteString(statusCode int, s string) error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to WriteString")
		return errors.New("superflous call to WriteString")
	}
	_, err := write(ctx, statusCode, []byte(s))
	if err != nil {
		return fmt.Errorf("error writing string to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// SendJSON writes a JSON object to the response body.
func (ctx *Context) JSON(statusCode int, a any) error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to JSON")
		return errors.New("superflous call to JSON")
	}
	ctx.Writer.Header().Set("Content-Type", "application/json")
	jsn, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = ctx.Writer.Write(jsn)
	if err != nil {
		return fmt.Errorf("error writing JSON to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// HTML writes an HTML template the response body.
func (ctx *Context) HTML(statusCode int, html string) error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to HTML")
		return errors.New("superflous call to HTML")
	}
	ctx.Writer.Header().Set("Content-Type", "text/html")
	_, err := write(ctx, statusCode, []byte(html))
	if err != nil {
		return fmt.Errorf("error writing HTML to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// Redirect redirects the request to a new URL.
func (ctx *Context) Redirect(path string) error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Redirect")
		return errors.New("superflous call to Redirect")
	}
	http.Redirect(ctx.Writer, ctx.request, path, http.StatusFound)
	ctx.isWritten = true
	return nil
}

// NotFound sets the HTTP status code.
func (ctx *Context) NotFound() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to NotFound")
		return errors.New("superflous call to NotFound")
	}
	_, err := write(ctx, http.StatusNotFound, []byte("Not Found"))
	if err != nil {
		return fmt.Errorf("error writing Not Found to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// MethodNotAllowed sets the HTTP status code 405 and writes a message to the response body.
func (ctx *Context) MethodNotAllowed() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to MethodNotAllowed")
		return errors.New("superflous call to MethodNotAllowed")
	}
	_, err := write(ctx, http.StatusMethodNotAllowed, []byte("Method Not Allowed"))
	if err != nil {
		return fmt.Errorf("error writing Method Not Allowed to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// BadRequest sets the HTTP status code 400 and writes a message to the response body.
func (ctx *Context) BadRequest() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to BadRequest")
		return errors.New("superflous call to BadRequest")
	}
	_, err := write(ctx, http.StatusBadRequest, []byte("Bad Request"))
	if err != nil {
		return fmt.Errorf("error writing Bad Request to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// InternalServerError sets the HTTP status code 500 and writes a message to the response body.
func (ctx *Context) InternalServerError() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to InternalServerError")
		return errors.New("superflous call to InternalServerError")
	}
	_, err := write(ctx, http.StatusInternalServerError, []byte("Internal Server Error"))
	if err != nil {
		return fmt.Errorf("error writing Internal Server Error to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// Unauthorized sets the HTTP status code 401 and writes a message to the response body.
func (ctx *Context) Unauthorized() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Unauthorized")
		return errors.New("superflous call to Unauthorized")
	}
	_, err := write(ctx, http.StatusUnauthorized, []byte("Unauthorized"))
	if err != nil {
		return fmt.Errorf("error writing Unauthorized to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// Forbidden sets the HTTP status code 403 and writes a message to the response body.
func (ctx *Context) Forbidden() error {
	if ctx.isWritten {
		ctx.Logger.Log(WARNING, "superflous call to Forbidden")
		return errors.New("superflous call to Forbidden")
	}
	_, err := write(ctx, http.StatusForbidden, []byte("Forbidden"))
	if err != nil {
		return fmt.Errorf("error writing Forbidden to response: %v", err)
	}
	ctx.isWritten = true
	return nil
}

// Request returns the HTTP request.
func (ctx *Context) Request() *http.Request {
	return ctx.request
}

// FormValue returns the value of the given form key.
func (ctx *Context) FormValue(key string) string {
	return ctx.request.FormValue(key)
}

// FormValues returns the values of the form.
func (ctx *Context) FormValues() map[string][]string {
	return ctx.request.Form
}

// extractParams extracts the parameters from the given path based on the provided pattern.
func extractParams(path, pattern string) map[string]string {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return nil
	}

	params := make(map[string]string, len(patternParts))
	for i := 1; i < len(patternParts); i++ {
		part := patternParts[i]
		if len(part) > 0 && part[0] == ':' {
			params[part[1:]] = pathParts[i]
		}
	}

	return params
}

// write writes data to the response body.
func write(ctx *Context, statusCode int, data []byte) (int, error) {
	ctx.Writer.WriteHeader(statusCode)
	return ctx.Writer.Write(data)
}
