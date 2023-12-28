package gort

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

	router.AddRoute("GET", "/users/:id", func(ctx *Context) {
		ctx.WriteString("hello")
	})

	route := router.Find("/users/foo")
	if route == nil {
		t.Error("unexpected nil route")
		return
	}

	if route.Pattern != "/users/:id" {
		t.Error("expected route pattern to be /users/:id")
		return
	}

	if route.Method != http.MethodGet {
		t.Error("expected route method to be GET")
		return
	}

	if route.Handler == nil {
		t.Error("unexpected nil handler")
		return
	}

	ctx := &Context{
		Params: extractParams("/users/foo", "/users/:id"),
	}

	if ctx.Params["id"] != "foo" {
		t.Error("expected params to be foo")
		return
	}

}

func TestWriteString(t *testing.T) {
	router := NewRouter()

	router.AddRoute("GET", "/users/:id", func(ctx *Context) {
		ctx.WriteString("hello")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/foo")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	if string(data) != "hello" {
		t.Error("expected response body to be hello")
		return
	}

}

func TestJSON(t *testing.T) {
	router := NewRouter()

	router.AddRoute("GET", "/json/test", func(ctx *Context) {
		ctx.JSON(map[string]string{
			"message": "hello",
		})
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/json/test")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	var jsnData map[string]string
	if err := json.Unmarshal(data, &jsnData); err != nil {
		t.Error(err)
		return
	}

	if jsnData["message"] != "hello" {
		t.Error("expected message to be hello")
		return
	}

}

func TestHTML(t *testing.T) {
	router := NewRouter()

	router.AddRoute("GET", "/html/test", func(ctx *Context) {
		ctx.HTML("<h1>hello</h1>")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/html/test")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, res.StatusCode)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	if string(data) != "<h1>hello</h1>" {
		t.Error("expected response body to be <h1>hello</h1>")
		return
	}
}

func TestRTree(t *testing.T) {
	tree := newRTree()

	tree.add(&Route{
		Method:  http.MethodGet,
		Pattern: "/users/:id",
		Handler: func(ctx *Context) {
			ctx.WriteString("hello")
		},
	})

	if tree.root == nil {
		t.Error("unexpected nil root")
		return
	}

	if tree.root.children == nil {
		t.Error("unexpected nil children")
		return
	}

	if len(tree.root.children) != 1 {
		t.Error("expected 1 child")
		return
	}

	if tree.root.children["users"] == nil {
		t.Error("unexpected nil child")
		return
	}

	route := tree.find("/users/foo")

	if route == nil {
		t.Error("unexpected nil route")
		return
	}
}
